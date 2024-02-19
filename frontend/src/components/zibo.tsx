import React, { useEffect, useState } from 'react';
import { Button, Card, Col, Descriptions, Divider, Row, Spin, Table, Tabs } from 'antd'
// import {FindZiboInstallationDetails} from "../../wailsjs/go/main/App";
import {
    BackupZiboInstallation,
    DownloadZibo,
    FindZiboInstallationDetails,
    GetAvailableLiveries,
    GetDownloadDetails,
    GetLiveries,
    InstallZibo,
    RestoreZiboInstallation,
    UpdateZibo
} from "../../wailsjs/go/main/App";
import { installer, utils } from "../../wailsjs/go/models";
import { BrowserOpenURL, LogInfo } from "../../wailsjs/runtime";
import InstalledLivery = installer.InstalledLivery;
import AvailableLivery = installer.AvailableLivery;

interface ZiboProps {
    installationDetails: utils.ZiboInstallation
}

function Zibo(props: ZiboProps) {
    // const screens = useBreakpoint();
    const [running, setRunning] = useState(false);
    const [progressDetails, setProgressDetails] = useState("")
    const [installedLiveries, setInstalledLiveries] = useState([] as InstalledLivery[]);
    const [availableLiveries, setAvailableLiveries] = useState([] as AvailableLivery[]);
    useEffect(() => {
        const fetchDetails = async () => {
            const details = await FindZiboInstallationDetails();
            const liveries = await GetLiveries(props.installationDetails);
            setInstalledLiveries(liveries);
        };
        // Call it once immediately
        fetchDetails();

    }, []);
    useEffect(() => {
        (async () => {
            const availableLiveries = await GetAvailableLiveries();
            LogInfo("availableLiveries:" + availableLiveries.length)
            setAvailableLiveries(availableLiveries)
        })()
    }, [])

    const handleBackup = async () => {
        setRunning(true);
        setProgressDetails("Backing up ...")
        await BackupZiboInstallation(props.installationDetails);
        const details = await FindZiboInstallationDetails();
        setRunning(false);
        window.location.reload();
    }

    const handleRestore = async () => {
        setRunning(true);
        setProgressDetails("Restoring ...")
        await RestoreZiboInstallation(props.installationDetails, "");
        const details = await FindZiboInstallationDetails();
        setRunning(false);
        window.location.reload();
    }
    const handleInstall = async () => {
        setRunning(true);
        const downloadInfo = await DownloadZibo(true);
        while (downloadInfo.isDownloading) {
            const downloadDetails = await GetDownloadDetails(false);
            setProgressDetails(`${downloadDetails.toFixed(2)}%`)
            if (downloadDetails === 100) {
                break;
            }
        }
        setProgressDetails("Installing ...")
        await InstallZibo(props.installationDetails, downloadInfo.path);
        const details = await FindZiboInstallationDetails();
        setRunning(false);
        window.location.reload();
    }

    const handleUpdate = async () => {
        setRunning(true);
        const downloadInfo = await DownloadZibo(false);
        while (downloadInfo.isDownloading) {
            const downloadDetails = await GetDownloadDetails(true);
            setProgressDetails(`${downloadDetails.toFixed(2)}%`)
            if (downloadDetails === 100) {
                break;
            }
        }
        setProgressDetails("Updating ...")
        const update_err = await UpdateZibo(props.installationDetails, downloadInfo.path);
        if (update_err != null) {
            setProgressDetails("Update failed, check log")
        } else {
            setProgressDetails("Update Succesfull")
        }
        const details = await FindZiboInstallationDetails();
        setRunning(false);
        window.location.reload();
    }

    return (
        <Card style={{
            minHeight: "100%",
        }}
        >
            <Spin spinning={running} tip={progressDetails}>
                <Row style={{ minHeight: "100%" }}>
                    <Col span={24} style={{
                        display: "flex",
                        justifyContent: "flex-end",
                        alignItems: "center",
                        height: "100%",
                        marginBottom: "12px"
                        // marginTop: "10%"
                        // marginRight: !screens.sm ? "12px" : "24px",
                        // marginLeft: !screens.sm ? "12px" : "24px",
                    }}>
                        <Button
                            type={"primary"}
                            danger={props.installationDetails.version !== ""}
                            style={{ marginRight: "12px" }}
                            onClick={handleInstall}>{props.installationDetails.version === "" ? "Install" : "Reinstall"}</Button>
                        <Button
                            type={"primary"}
                            style={{ marginRight: "12px" }}
                            disabled={
                                props.installationDetails.version === props.installationDetails.remoteVersion ||
                                props.installationDetails.version === ""
                            }
                            onClick={handleUpdate}>Update</Button>
                        <Button
                            type={"primary"}
                            style={{ marginRight: "12px" }}
                            disabled={props.installationDetails.version === ""}
                            onClick={handleBackup}>Backup</Button>
                        <Button
                            type={"primary"}
                            danger={true}
                            style={{ marginRight: "12px" }}
                            disabled={props.installationDetails.backupVersion === "N/A"}
                            onClick={handleRestore}>Restore</Button>
                    </Col>
                    <Col span={24}>
                        <Descriptions
                            // title="Installation Details"
                            layout="vertical"
                            bordered
                            // column={2}
                            items={[
                                {
                                    key: '1',
                                    label: 'Installed Version',
                                    children: `${props.installationDetails.version}`,

                                },
                                {
                                    key: '2',
                                    label: 'Current Version',
                                    children: `${props.installationDetails.remoteVersion}`,
                                },
                                {
                                    key: '3',
                                    label: 'Backup Version',
                                    children: `${props.installationDetails.backupVersion}`,
                                },
                                {
                                    key: '10',
                                    label: 'Installed Path',
                                    children: `${props.installationDetails.path}`,
                                },

                            ]}
                        />
                    </Col>
                    <Divider />
                    <Col span={24}>
                        <Tabs items={[
                            {
                                label: "Installed Liveries",
                                key: "installedLiveries",
                                children: <Table
                                    // title={() => "Liveries"}
                                    dataSource={installedLiveries}
                                    style={{ overflow: "scroll" }}
                                    columns={[
                                        {
                                            title: 'Icon',
                                            dataIndex: 'icon',
                                            key: 'icon',
                                            render: (icon: string) => <img src={`data:image/png;base64,${icon}`}
                                                alt={"icon"}
                                                width={160} height={90} />
                                        },
                                        {
                                            title: 'Name',
                                            dataIndex: 'name',
                                            key: 'name',
                                        },
                                        // {
                                        //     title: 'Installed Path',
                                        //     dataIndex: 'path',
                                        //     key: 'path',
                                        // },
                                    ]}
                                />,
                            },
                            {
                                label: "Available Liveries",
                                key: "availableLiveries",
                                children: <Table
                                    title={() => "Liveries"}
                                    dataSource={availableLiveries}
                                    style={{ overflow: "scroll" }}
                                    pagination={{ pageSize: 25 }}
                                    columns={[
                                        {
                                            title: 'Icon',
                                            dataIndex: 'icon',
                                            key: 'icon',
                                            render: (icon: string) => <img src={`data:image/png;base64,${icon}`}
                                                alt={"icon"}
                                                width={320} height={200} />
                                        },
                                        {
                                            title: 'Name',
                                            dataIndex: 'name',
                                            key: 'name',
                                        },
                                        {
                                            title: '',
                                            dataIndex: 'url',
                                            key: 'url',
                                            render: (url: string) => <Button
                                                onClick={() => {
                                                    BrowserOpenURL(url)
                                                }} title={'DOWNLOAD'}>DOWNLOAD</Button>
                                        },
                                    ]}
                                />,
                            },
                        ]}>

                        </Tabs>

                    </Col>

                </Row>
            </Spin>
        </Card>
    )
}

export default Zibo
