// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {utils} from '../models';
import {main} from '../models';
import {installer} from '../models';

export function BackupZiboInstallation(arg1:utils.ZiboInstallation):Promise<boolean>;

export function CheckXPlanePath(arg1:string,arg2:Array<string>):Promise<boolean>;

export function DeleteFiles(arg1:Array<string>):Promise<string>;

export function DownloadZibo(arg1:boolean):Promise<main.DownloadInfo>;

export function FindZiboInstallationDetails():Promise<Array<utils.ZiboInstallation>>;

export function GetAvailableLiveries():Promise<Array<installer.AvailableLivery>>;

export function GetBackups():Promise<Array<installer.ZiboBackup>>;

export function GetCachedFiles():Promise<Array<utils.CachedFile>>;

export function GetConfig():Promise<utils.Config>;

export function GetDownloadDetails(arg1:boolean):Promise<number>;

export function GetLiveries(arg1:utils.ZiboInstallation):Promise<Array<installer.InstalledLivery>>;

export function GetOs():Promise<string>;

export function InstallZibo(arg1:utils.ZiboInstallation,arg2:string):Promise<void>;

export function IsXPlanePathConfigured():Promise<boolean>;

export function OpenDirDialog():Promise<string>;

export function RestoreZiboInstallation(arg1:utils.ZiboInstallation,arg2:string):Promise<boolean>;

export function UpdateZibo(arg1:utils.ZiboInstallation,arg2:string):Promise<void>;
