export namespace installer {
	
	export class AvailableLivery {
	    name: string;
	    url: string;
	    source: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new AvailableLivery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.source = source["source"];
	        this.icon = source["icon"];
	    }
	}
	export class InstalledLivery {
	    name: string;
	    path: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new InstalledLivery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.icon = source["icon"];
	    }
	}
	export class ZiboBackup {
	    backupPath: string;
	    version: string;
	    date: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ZiboBackup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.backupPath = source["backupPath"];
	        this.version = source["version"];
	        this.date = source["date"];
	        this.size = source["size"];
	    }
	}

}

export namespace main {
	
	export class DownloadInfo {
	    isDownloading: boolean;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isDownloading = source["isDownloading"];
	        this.path = source["path"];
	    }
	}

}

export namespace utils {
	
	export class CachedFile {
	    name: string;
	    path: string;
	    completedSize: number;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new CachedFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.completedSize = source["completedSize"];
	        this.size = source["size"];
	    }
	}
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

