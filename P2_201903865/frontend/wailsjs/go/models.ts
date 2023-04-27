export namespace sys {
	
	export class CPUUsage {
	    avg: number;
	
	    static createFrom(source: any = {}) {
	        return new CPUUsage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avg = source["avg"];
	    }
	}
	export class DISKUsage {
	    used: number;
	    free: number;
	
	    static createFrom(source: any = {}) {
	        return new DISKUsage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.used = source["used"];
	        this.free = source["free"];
	    }
	}
	export class RAMUsage {
	    used: number;
	    free: number;
	
	    static createFrom(source: any = {}) {
	        return new RAMUsage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.used = source["used"];
	        this.free = source["free"];
	    }
	}
	export class USBport {
	    vid: string;
	    pid: string;
	    name: string;
	    device: string;
	    status: boolean;
	    port_name: string;
	
	    static createFrom(source: any = {}) {
	        return new USBport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vid = source["vid"];
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.device = source["device"];
	        this.status = source["status"];
	        this.port_name = source["port_name"];
	    }
	}

}

