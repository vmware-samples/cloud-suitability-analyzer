/*******************************************************************************
 * Copyright (c) 2018 - Present VMWare, Inc. All Rights Reserved.
 ******************************************************************************/

package model

//Created By BootstrapBinsTemplate.txt found under go/resources folder
//Created @ 2021-01-15 14:53:04.807469154 +0000 UTC m=+0.087387403

func BootstrapBins() []Bin {
    var BootstrapBins = []Bin{
        
            { Name: "TKG",
            Tags:
            []*BinTag{  { Name: "docker", Type: 1, Action: "OR", },  { Name: "stateful", Type: 0, Action: "AND", },  { Name: "javaee", Type: 0, Action: "AND", },  { Name: "fullprofile", Type: 0, Action: "AND", },  { Name: "jni", Type: 1, Action: "OR", },  { Name: "nonstandard-protocol", Type: 1, Action: "OR", },  { Name: "corba", Type: 1, Action: "OR", },  },
             },
        
            { Name: "TAS",
            Tags:
            []*BinTag{  { Name: "webprofile", Type: 1, Action: "OR", },  { Name: "spring", Type: 1, Action: "OR", },  { Name: "spring-boot", Type: 1, Action: "OR", },  { Name: "webcontainer", Type: 1, Action: "OR", },  { Name: "rest", Type: 1, Action: "OR", },  },
             },
        
            { Name: "BATCH",
            Tags:
            []*BinTag{  { Name: "batch", Type: 1, Action: "OR", },  { Name: "etl", Type: 1, Action: "OR", },  },
             },
        
            { Name: "CONTAINER",
            Tags:
            []*BinTag{  { Name: "term", Type: 1, Action: "OR", },  { Name: "metrics", Type: 1, Action: "OR", },  { Name: "sudo", Type: 1, Action: "OR", },  { Name: "healthcheck", Type: 1, Action: "OR", },  { Name: "ehcache", Type: 1, Action: "OR", },  { Name: "non-root", Type: 1, Action: "OR", },  { Name: "hardip", Type: 1, Action: "OR", },  { Name: "processexit", Type: 1, Action: "OR", },  { Name: "distcache", Type: 1, Action: "OR", },  { Name: "wscluster", Type: 1, Action: "OR", },  { Name: "wlcluster", Type: 0, Action: "ORs", },  { Name: "io", Type: 1, Action: "OR", },  { Name: "log2file", Type: 1, Action: "OR", },  { Name: "docker", Type: 1, Action: "OR", },  { Name: "transaction", Type: 1, Action: "OR", },  },
             },
        
            { Name: "BOOT",
            Tags:
            []*BinTag{  { Name: "cdi", Type: 1, Action: "OR", },  { Name: "ejb", Type: 1, Action: "OR", },  { Name: "jaxws", Type: 1, Action: "OR", },  { Name: "jsf", Type: 1, Action: "OR", },  { Name: "mdb", Type: 1, Action: "OR", },  { Name: "struts", Type: 1, Action: "OR", },  { Name: "txn", Type: 1, Action: "OR", },  { Name: "webservlet", Type: 1, Action: "OR", },  },
             },
        
            { Name: "Modern Microservices",
            Tags:
            []*BinTag{  { Name: "spring-boot", Type: 0, Action: "AND", },  { Name: "jar", Type: 0, Action: "AND", },  },
             },
        
            { Name: "Spring",
            Tags:
            []*BinTag{  { Name: "spring", Type: 1, Action: "OR", },  },
             },
        
            { Name: "Stateful",
            Tags:
            []*BinTag{  { Name: "stateful", Type: 1, Action: "OR", },  { Name: "spring", Type: 2, Action: "EXCLUDE", },  },
             },
        
            { Name: "Web MVC",
            Tags:
            []*BinTag{  { Name: "ui", Type: 1, Action: "OR", },  { Name: "webapp", Type: 0, Action: "AND", },  { Name: "webcontainer", Type: 0, Action: "AND", },  { Name: "jsp", Type: 1, Action: "OR", },  { Name: "jsf", Type: 1, Action: "OR", },  { Name: "portlet", Type: 1, Action: "OR", },  { Name: "struts", Type: 1, Action: "OR", },  { Name: "jetty", Type: 1, Action: "OR", },  },
             },
        
            { Name: "Thick Java Clients",
            Tags:
            []*BinTag{  { Name: "fat-client", Type: 1, Action: "OR", },  },
             },
        
            { Name: "Caching",
            Tags:
            []*BinTag{  { Name: "cache", Type: 1, Action: "OR", },  { Name: "persistence", Type: 1, Action: "OR", },  },
             },
        
            { Name: "Data",
            Tags:
            []*BinTag{  { Name: "io", Type: 1, Action: "OR", },  },
             },
        
            { Name: "Monolith",
            Tags:
            []*BinTag{  { Name: "soap", Type: 0, Action: "AND", },  { Name: "webservices", Type: 0, Action: "AND", },  { Name: "messaging", Type: 0, Action: "AND", },  { Name: "ear", Type: 0, Action: "AND", },  { Name: "stateful", Type: 0, Action: "AND", },  { Name: "rpc", Type: 1, Action: "OR", },  { Name: "ant", Type: 1, Action: "OR", },  { Name: "websphere", Type: 1, Action: "OR", },  { Name: "weblogic", Type: 1, Action: "OR", },  { Name: "glassfish", Type: 1, Action: "OR", },  { Name: "jboss", Type: 1, Action: "OR", },  { Name: "security", Type: 1, Action: "OR", },  },
             },
        
            { Name: "FAAS",
            Tags:
            []*BinTag{  { Name: "persistence", Type: 1, Action: "OR", },  { Name: "security", Type: 1, Action: "OR", },  { Name: "authentication", Type: 1, Action: "OR", },  { Name: "boottime", Type: 1, Action: "OR", },  },
             },
        
    }
    return BootstrapBins
}