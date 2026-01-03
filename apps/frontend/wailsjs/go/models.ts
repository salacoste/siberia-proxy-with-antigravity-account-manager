export namespace accounts {
	
	export class AccountDTO {
	    id: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    email: string;
	    recovery_email: string;
	    proxy_group: string;
	    is_active: boolean;
	    stats: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.email = source["email"];
	        this.recovery_email = source["recovery_email"];
	        this.proxy_group = source["proxy_group"];
	        this.is_active = source["is_active"];
	        this.stats = source["stats"];
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

}

export namespace config {
	
	export class AppConfig {
	    theme: string;
	    language: string;
	    auto_refresh: boolean;
	    refresh_interval: number;
	    db_sync: boolean;
	    proxy_port: number;
	    upstream_proxy: string;
	    zai_enabled: boolean;
	    zai_base_url: string;
	    zai_api_key: string;
	    mitm_enabled: boolean;
	    auth_enabled: boolean;
	    auth_token: string;
	    master_key: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.auto_refresh = source["auto_refresh"];
	        this.refresh_interval = source["refresh_interval"];
	        this.db_sync = source["db_sync"];
	        this.proxy_port = source["proxy_port"];
	        this.upstream_proxy = source["upstream_proxy"];
	        this.zai_enabled = source["zai_enabled"];
	        this.zai_base_url = source["zai_base_url"];
	        this.zai_api_key = source["zai_api_key"];
	        this.mitm_enabled = source["mitm_enabled"];
	        this.auth_enabled = source["auth_enabled"];
	        this.auth_token = source["auth_token"];
	        this.master_key = source["master_key"];
	    }
	}

}

export namespace updater {
	
	export class UpdateInfo {
	    available: boolean;
	    current_version: string;
	    latest_version: string;
	    download_url: string;
	    release_notes: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.current_version = source["current_version"];
	        this.latest_version = source["latest_version"];
	        this.download_url = source["download_url"];
	        this.release_notes = source["release_notes"];
	        this.error = source["error"];
	    }
	}

}

