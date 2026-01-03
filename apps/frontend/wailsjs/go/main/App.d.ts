export function ListAccounts(): Promise<Array<any>>;
export function DeleteAccount(arg1: number): Promise<void>;
export function GetAppConfig(): Promise<any>;
export function UpdateAppConfig(arg1: any): Promise<void>;
export function CreateAccount(email: string, password: string, recovery: string, proxyGroup: string): Promise<void>;
export function ActivateAccount(id: number): Promise<void>;
export function Greet(arg1: string): Promise<string>;
