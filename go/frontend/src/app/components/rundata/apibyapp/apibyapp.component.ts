import { Component, OnInit } from '@angular/core';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';
import {RundataService} from '../../../services/rundata.service';
import {Apiusage} from '../../../model/appapiusage';
import {
  ClrDatagridFilterInterface,
  ClrDatagridNumericFilterInterface,
  ClrDatagridStringFilterInterface
} from '@clr/angular';
import {Subject} from "rxjs";

@Component({
  selector: 'apibyapp',
  templateUrl: './apibyapp.component.html',
  styleUrls: ['./apibyapp.component.css', '../rundatasummary/rundatasummary.component.css']
})
export class ApiByAppComponent implements OnInit {

  public searchCrit: any = '';
  public fileName: string;

  public gridData: (string|number)[][] = [];
  public gridColumns: string[] = [];

  public appNameFilter = new AppNameFilter();

  constructor(private router: Router, private route: ActivatedRoute, private rundataService: RundataService) {

  }

  ngOnInit(): void {
    this.route.paramMap.subscribe((params: ParamMap) => {
      this.resetPage();
      const runId = Number(params.get('id'));
      console.log('runid is : ' + runId);
      this.fetchApiByApps(runId);
    });
  }

  resetPage(): void {
    this.gridData = [];
    this.gridColumns = [];
  }

  fetchApiByApps(runId: number): void {
    this.rundataService.getApiByAppUsage(runId).subscribe(appApiUsageReturned => {
      appApiUsageReturned.cols.forEach(col => {
        if (col !== 'App') {
          this.gridColumns.push(col);
        }
      });
      if (appApiUsageReturned.cols && appApiUsageReturned.data) {
        console.log(this.gridData.length);
        appApiUsageReturned.data.forEach(datum => {
          const innerArr: (string|number)[] = [];
          innerArr.push(datum.application);
          appApiUsageReturned.cols.forEach(col => {
            let apiusage: Apiusage;
            apiusage = datum.apiusage.find(indApiUsage => indApiUsage.api.toLowerCase() == col.toLowerCase());
            if (apiusage) {
              innerArr.push(apiusage.usageCount);
            }
          });
          this.gridData.push(innerArr);
        });
        console.log(this.gridColumns);
      }
    }, error => {
      console.log(error);
    });
  }

}

class AppNameFilter implements ClrDatagridStringFilterInterface<[]> {
  accepts(apiByAppRow: [], search: string): boolean {
    // @ts-ignore
    return apiByAppRow[0].toLowerCase().indexOf(search) >= 0;
  }
}