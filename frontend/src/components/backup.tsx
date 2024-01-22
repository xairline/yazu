import React, {useEffect, useState} from 'react';
import {
    DeleteFiles,
    FindZiboInstallationDetails,
    GetBackups,
    GetCachedFiles,
    RestoreZiboInstallation
} from "../../wailsjs/go/main/App";
import {Button, Card, Dropdown, Row, Spin, Table} from "antd";
import {installer, utils} from "../../wailsjs/go/models";
import ZiboBackup = installer.ZiboBackup;
import ZiboInstallation = utils.ZiboInstallation;


function Backup() {
    const [backups, setBackups] = useState([] as any);
    const [running, setRunning] = useState(false);
    const [cachedFiles, setCachedFiles] = useState({});
    const [progressDetails, setProgressDetails] = useState("")
    const [ziboDetails, setZiboDetails] = useState({} as ZiboInstallation[]);
    const [selectedBackups, setSelectedBackups] = useState([] as any);
    const [trigger, setTrigger] = useState(0);
    useEffect(() => {
        (async () => {
            const backups = await GetBackups();
            const cachedFiles = await GetCachedFiles();
            const details = await FindZiboInstallationDetails();
            setZiboDetails(details)
            setBackups(backups)
            setCachedFiles(cachedFiles)
        })();

    }, [trigger]);


    return (

        <Spin spinning={running} tip={progressDetails}>
            <Row>
                <Button
                    danger type={"primary"}
                    disabled={selectedBackups.length <= 0}
                    onClick={
                        async () => {
                            setRunning(true)
                            setProgressDetails("Deleting ...")
                            await DeleteFiles(selectedBackups.map((backup: ZiboBackup) => {
                                return backup.backupPath
                            }))
                            setRunning(false)
                            setTrigger(prev => prev + 1)
                        }
                    }
                > Delete </Button>
            </Row>
            <Card style={{
                minHeight: "100%",
                marginTop: "12px",
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
                            key: backup.backupPath,
                            version: backup.version,
                            date: new Date(date1 + "T" + date2 + "Z"),
                            size: (backup.size / 1024 / 1024).toFixed(2) + "MB",
                            backupPath: backup.backupPath,
                        }
                    })}
                    rowSelection={{
                        type: "checkbox",
                        onChange: (selectedRowKeys: any, selectedRows: any) => {
                            setSelectedBackups(selectedRows)
                        }
                    }}
                    columns={[
                        {
                            title: 'Version',
                            dataIndex: 'version',
                            key: 'key',
                            sorter: (a: any, b: any) => {
                                const parseSemver = (version: string) => {
                                    return version.split('.').map((num: string) => parseInt(num, 10));
                                };
                                const [majorA, minorA, patchA] = parseSemver(a.version);
                                const [majorB, minorB, patchB] = parseSemver(b.version);

                                if (majorA !== majorB) {
                                    return majorA - majorB;
                                }
                                if (minorA !== minorB) {
                                    return minorA - minorB;
                                }
                                return patchA - patchB;
                            },
                            defaultSortOrder: 'descend',
                        },
                        {
                            title: 'Date',
                            dataIndex: 'date',
                            key: 'date',
                            render: (date: string) => {
                                return date.toLocaleString()
                            },
                            sorter: (a: any, b: any) => {
                                // Convert date strings to Date objects
                                const dateA = new Date(a.date);
                                const dateB = new Date(b.date);

                                // Compare the dates
                                // @ts-ignore
                                return dateA - dateB;
                            },
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
