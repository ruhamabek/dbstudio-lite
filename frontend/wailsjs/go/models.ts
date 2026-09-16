export namespace db {
	
	export class Column {
	    Name: string;
	    DataType: string;
	    IsNullable: boolean;
	    DefaultVal?: string;
	
	    static createFrom(source: any = {}) {
	        return new Column(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.DataType = source["DataType"];
	        this.IsNullable = source["IsNullable"];
	        this.DefaultVal = source["DefaultVal"];
	    }
	}
	export class Config {
	    driver: string;
	    host: string;
	    port: number;
	    user: string;
	    password: string;
	    database: string;
	    sslMode: string;
	    filePath: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.driver = source["driver"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.database = source["database"];
	        this.sslMode = source["sslMode"];
	        this.filePath = source["filePath"];
	    }
	}
	export class QueryResult {
	    columns: string[];
	    rows: any[][];
	    rowsAffected: number;
	    durationMs: number;
	
	    static createFrom(source: any = {}) {
	        return new QueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.columns = source["columns"];
	        this.rows = source["rows"];
	        this.rowsAffected = source["rowsAffected"];
	        this.durationMs = source["durationMs"];
	    }
	}
	export class SavedConnection {
	    id: string;
	    name: string;
	    url: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new SavedConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.url = source["url"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
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
	export class Table {
	    Schema: string;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new Table(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Schema = source["Schema"];
	        this.Name = source["Name"];
	    }
	}

}

