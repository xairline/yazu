export namespace utils {
	
	export class Config {
	    XPlanePath: string;
	    YazuCachePath: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.XPlanePath = source["XPlanePath"];
	        this.YazuCachePath = source["YazuCachePath"];
	    }
	}
	export class ZiboInstallation {
	    path: string;
	    version: string;
	    remoteVersion: string;
	    backupVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new ZiboInstallation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.version = source["version"];
	        this.remoteVersion = source["remoteVersion"];
	        this.backupVersion = source["backupVersion"];
	    }
	}

}

