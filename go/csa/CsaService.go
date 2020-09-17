/*******************************************************************************
 * Copyright (c) 2018 - Present VMWare, Inc. All Rights Reserved.
 * SPDX-License-Identifier: BSD-2
 ******************************************************************************/

package csa

import (
	"fmt"
	"os"
	"sync"

	"github.com/antchfx/xmlquery"
	"github.com/vmware-samples/cloud-suitability-analyzer/go/db"
	"github.com/vmware-samples/cloud-suitability-analyzer/go/model"
	"github.com/vmware-samples/cloud-suitability-analyzer/go/report"
	"github.com/vmware-samples/cloud-suitability-analyzer/go/util"
	"gopkg.in/yaml.v3"
)

//The Engine that does file parsing and rule matching
type CsaService struct {
	ruleRepository       db.RuleRepository
	runRepository        db.RunRepository
	findingRepository    db.FindingRepository
	reportDataRepository db.ReportDataRepository
	slocRepository       db.SlocRepository
	scoringRepository    db.ScoringRepository
	reportService        *report.ReportService
	fileUtil             *util.FileUtil
	saveChan             chan interface{} // = make(chan interface{}, *util.MaxBuffer)
	indexChan            chan interface{}
	saveDone             chan interface{} //  = make(chan interface{}, *util.MaxBuffer)
	indexDone            chan interface{} //  = make(chan interface{}, *util.MaxBuffer)
	analysisDone         bool             //= false
	savingDone           bool             //= false
	findingsSaved        int              //= 0
	findingsIndexed      int              //= 0
	yamlDocs             map[string](*yaml.Node)
	yamlMux              sync.Mutex
	xmlDocs              map[string](*xmlquery.Node)
	xmlMux               sync.Mutex
}

func NewCsaSvc(mgr *db.Repositories) *CsaService {
	return NewCsaService(mgr.Rules, mgr.Run, mgr.Findings, mgr.Reports, mgr.Sloc, mgr.Scoring, report.NewReportSvc(mgr))
}

func NewCsaService(ruleRepository db.RuleRepository,
	runRepository db.RunRepository,
	findingRepository db.FindingRepository,
	reportDataRepository db.ReportDataRepository,
	slocRepository db.SlocRepository,
	scoringRepo db.ScoringRepository,
	reportService *report.ReportService) *CsaService {

	return &CsaService{
		ruleRepository:       ruleRepository,
		runRepository:        runRepository,
		findingRepository:    findingRepository,
		reportDataRepository: reportDataRepository,
		slocRepository:       slocRepository,
		scoringRepository:    scoringRepo,
		reportService:        reportService,
		fileUtil:             util.NewFileUtil(),
		saveChan:             make(chan interface{}, *util.MaxBuffer),
		indexChan:            make(chan interface{}, *util.MaxBuffer),
		saveDone:             make(chan interface{}),
		indexDone:            make(chan interface{}),
		analysisDone:         false,
		savingDone:           false,
		findingsSaved:        0,
		findingsIndexed:      0,
		xmlDocs:              make(map[string](*xmlquery.Node)),
		yamlDocs:             make(map[string](*yaml.Node)),
	}

}

func (csaService *CsaService) PerformAnalysis(run *model.Run) {

	csaService.startRun(run)
	csaService.gatherFiles(run)
	if !util.ProcessHadErrors("gathering") {
		if !*util.WriteConfigsOnly {
			if len(run.Applications) > 0 {
				saveWorkerCnt, indexWorkerCnt := csaService.startWorkers(run)
				if *util.SerialAppAnalysis {
					csaService.SerialAnalysis(run)
				} else {
					csaService.concurrentAnalysis(run)
				}
				csaService.waitForSavingAndIndexingToComplete(run, saveWorkerCnt, indexWorkerCnt)
				csaService.generateSloc(run)
				csaService.scoreApps(run)
				csaService.generateReports(run)
			}
		}
	}

	csaService.stopRun(run)

	if util.HasErrors() && *util.DisplayErrors {
		util.DumpErrors(run.ID)
	} else if !*util.WriteConfigsOnly {
		//Generate these last so they don't distort run duration measurement! (even though its only milliseconds...)
		csaService.genRuleMetrics(run)
	}

}

func (csaService *CsaService) concurrentAnalysis(run *model.Run) {

	var errors []error

	run.StartActivity("analysis")
	waitGroup := sync.WaitGroup{}

	apps := run.AppsOrdered()
	util.InitializeSpinners(len(apps))

	//app analysis
	for i := range apps {
		waitGroup.Add(1)
		go func(idx int) {
			defer waitGroup.Done()
			errors = append(errors, csaService.analyzeApp(run, apps[idx], csaService.saveChan)...)
		}(i)
	}

	waitGroup.Wait()

	msg := "done!"

	if util.ProcessHadErrors("Analysis") {
		msg = "errors!"
	}

	run.StopActivityLF("analysis", fmt.Sprintf("Analyzing...%s", msg), false, true)

	csaService.analysisDone = true
	close(csaService.saveChan)

}

func (csaService *CsaService) SerialAnalysis(run *model.Run) {
	var errors []error

	run.StartActivity("analysis")

	apps := run.AppsOrdered()
	util.InitializeSpinners(1)

	//app analysis
	for i := range apps {
		errors = append(errors, csaService.analyzeApp(run, apps[i], csaService.saveChan)...)
	}

	msg := "done!"

	if util.ProcessHadErrors("Analysis") {
		msg = "errors!"
	}

	run.StopActivityLF("analysis", fmt.Sprintf("Analyzing...%s", msg), false, true)

	csaService.analysisDone = true
	close(csaService.saveChan)
}

func (csaService *CsaService) startRun(run *model.Run) {
	err := csaService.runRepository.StartRun(run)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Starting Analysis Run! Details: %v", err)
		os.Exit(1)
	}
}

func (csaService *CsaService) stopRun(run *model.Run) {
	err := csaService.runRepository.StopRun(run)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Stoping Analysis Run! Details: %v", err)
	}
}
