import React, {useEffect, useState} from 'react';
import {
    FindZiboInstallationDetails,
    GetBackups,
    GetCachedFiles,
    RestoreZiboInstallation
} from "../../wailsjs/go/main/App";
import {Card, Dropdown, Spin, Table} from "antd";
import {installer, utils} from "../../wailsjs/go/models";
import ZiboBackup = installer.ZiboBackup;
import ZiboInstallation = utils.ZiboInstallation;


function Backup() {
    const [backups, setBackups] = useState([] as any);
    const [running, setRunning] = useState(false);
    const [cachedFiles, setCachedFiles] = useState({});
    const [progressDetails, setProgressDetails] = useState("")
    const [ziboDetails, setZiboDetails] = useState({} as ZiboInstallation[]);
    useEffect(() => {
        (async () => {
            const backups = await GetBackups();
            const cachedFiles = await GetCachedFiles();
            const details = await FindZiboInstallationDetails();
            setZiboDetails(details)
            setBackups(backups)
            setCachedFiles(cachedFiles)
        })();

    }, []);


    return (

        <Spin spinning={running} tip={progressDetails}>
            <Card style={{
                minHeight: "100%",
            }}
            >
                <Table
                    dataSource={backups?.map((backup: ZiboBackup) => {
                        const date = backup.date
                        let date1 = date.split("_")[0]
                        let date2 = date.split("_")[1]
                        date2 = date2.replace(/-/g, ":")
                        // // return date1 + "T" + date2 + "Z"
                        // return new Date(date1 + "T" + date2 + "Z").toLocaleString()
                        return {
                            key: backup.version,
                            version: backup.version,
                            date: new Date(date1 + "T" + date2 + "Z"),
                            size: (backup.size / 1024 / 1024).toFixed(2) + "MB",
                            backupPath: backup.backupPath,
                        }
                    })}
                    columns={[
                        {
                            title: 'Version',
                            dataIndex: 'version',
                            key: 'version',
                            sorter: (a: any, b: any) => {
                                return parseInt(a.version.split(".")[2]) - parseInt(b.version.split(".")[2])
                            },
                        },
                        {
                            title: 'Date',
                            dataIndex: 'date',
                            key: 'date',
                            render: (date: string) => {
                                return date.toLocaleString()
                            },
                            sorter: (a: any, b: any) => a.date - b.date,
                        },
                        {
                            title: 'Size',
                            dataIndex: 'size',
                            key: 'size',
                        },
                        {
                            title: 'Action',
                            dataIndex: 'action',
                            key: 'action',
                            render: (text: string, record: ZiboBackup) => {
                                // return <Button danger={true} type={"primary"} onClick={
                                //     async () => {
                                //         setRunning(true)
                                //         setProgressDetails("Restoring ...")
                                //         await RestoreZiboInstallation(ziboDetails[0], record.backupPath)
                                //         setRunning(false)
                                //     }
                                // }>Restore</Button>
                                return <Dropdown.Button
                                    danger={true}
                                    type={"primary"}
                                    menu={{
                                        items: ziboDetails.map((ziboDetail: ZiboInstallation) => {
                                            return {
                                                label: ziboDetail.path.split("/Aircraft/")[1].split("/plugins/")[0],
                                                key: ziboDetail.path,
                                            }
                                        }),
                                        onClick: async (e) => {
                                            setRunning(true)
                                            setProgressDetails(`Restoring ...${e.key}`)
                                            const installationDetails = ziboDetails.find(
                                                (ziboDetail: ZiboInstallation) => {
                                                    return ziboDetail.path === e.key
                                                }
                                            )
                                            if (!installationDetails) {
                                                return
                                            }
                                            await RestoreZiboInstallation(installationDetails, record.backupPath)
                                            setRunning(false)
                                        }
                                    }}
                                >
                                    Restore
                                </Dropdown.Button>
                            }
                        },
                    ]}
                />
            </Card>
            {/*<Divider/>*/}
            {/*<Card style={{*/}
            {/*    minHeight: "100%",*/}
            {/*}}*/}
            {/*>*/}
            {/*    <Skeleton/>*/}
            {/*    {JSON.stringify(cachedFiles)}*/}
            {/*</Card>*/}
        </Spin>
    )
}

export default Backup
