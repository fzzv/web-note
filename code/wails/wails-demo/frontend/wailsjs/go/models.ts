export namespace main {
	
	export class DownloadState {
	    status: string;
	    message: string;
	    url: string;
	    fileName: string;
	    destination: string;
	    downloadedBytes: number;
	    totalBytes: number;
	    progress: number;
	
	    static createFrom(source: any = {}) {
	        return new DownloadState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.message = source["message"];
	        this.url = source["url"];
	        this.fileName = source["fileName"];
	        this.destination = source["destination"];
	        this.downloadedBytes = source["downloadedBytes"];
	        this.totalBytes = source["totalBytes"];
	        this.progress = source["progress"];
	    }
	}

}

