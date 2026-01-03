export namespace config {
    export class AppConfig {
        theme: string;
        language: string;
        auto_refresh: boolean;
        refresh_interval: number;
        db_sync: boolean;
        proxy_port: number;
        upstream_proxy: string;

        // Zai
        zai_enabled: boolean;
        zai_base_url: string;
        zai_api_key: string;

        // Auth
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
            this.auth_enabled = source["auth_enabled"];
            this.auth_token = source["auth_token"];
            this.master_key = source["master_key"];
        }
    }
}
