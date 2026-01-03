export namespace updater {

    export class UpdateInfo {
        available: boolean;
        version: string;
        release_url: string;
        download_url: string;
        description: string;

        static createFrom(source: any = {}) {
            return new UpdateInfo(source);
        }

        constructor(source: any = {}) {
            if ('string' === typeof source) source = JSON.parse(source);
            this.available = source["available"];
            this.version = source["version"];
            this.release_url = source["release_url"];
            this.download_url = source["download_url"];
            this.description = source["description"];
        }
    }

}
