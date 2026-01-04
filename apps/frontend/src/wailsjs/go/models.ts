export namespace config {

	export class AppConfig {
		theme: string;
		language: string;
		auto_refresh: boolean;
		refresh_interval: number;
		db_sync: boolean;

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
			this.target_ide = source["target_ide"];
		}
		target_ide: string;
	}

}

export * from './updater/models';
