export namespace adb {
	
	export class DeviceInfo {
	    serial: string;
	    state: string;
	    model: string;
	    product: string;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serial = source["serial"];
	        this.state = source["state"];
	        this.model = source["model"];
	        this.product = source["product"];
	    }
	}

}

export namespace config {
	
	export class AppConfig {
	    sdkDir: string;
	    theme: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sdkDir = source["sdkDir"];
	        this.theme = source["theme"];
	    }
	}
	export class Phone {
	    name: string;
	    serial: string;
	    apkPath: string;
	    reinstall: boolean;
	    allowDowngrade: boolean;
	    grantPermissions: boolean;
	    autoLaunch: boolean;
	    launchPackage: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new Phone(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.serial = source["serial"];
	        this.apkPath = source["apkPath"];
	        this.reinstall = source["reinstall"];
	        this.allowDowngrade = source["allowDowngrade"];
	        this.grantPermissions = source["grantPermissions"];
	        this.autoLaunch = source["autoLaunch"];
	        this.launchPackage = source["launchPackage"];
	        this.notes = source["notes"];
	    }
	}

}

export namespace main {
	
	export class SDKStatus {
	    sdkDir: string;
	    adbPath: string;
	    version: string;
	    ready: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SDKStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sdkDir = source["sdkDir"];
	        this.adbPath = source["adbPath"];
	        this.version = source["version"];
	        this.ready = source["ready"];
	    }
	}
	export class AppState {
	    config: config.AppConfig;
	    sdk: SDKStatus;
	
	    static createFrom(source: any = {}) {
	        return new AppState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config = this.convertValues(source["config"], config.AppConfig);
	        this.sdk = this.convertValues(source["sdk"], SDKStatus);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CheckResult {
	    ok: boolean;
	    message: string;
	    device?: adb.DeviceInfo;
	
	    static createFrom(source: any = {}) {
	        return new CheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.message = source["message"];
	        this.device = this.convertValues(source["device"], adb.DeviceInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InstallOutcome {
	    success: boolean;
	    message: string;
	    elapsedMs: number;
	
	    static createFrom(source: any = {}) {
	        return new InstallOutcome(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.elapsedMs = source["elapsedMs"];
	    }
	}

}

