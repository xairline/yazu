import React, {useEffect, useState} from 'react';
import {Button, Card, Col, Descriptions, Divider, Row, Spin, Table, Tabs} from 'antd'
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
import {installer, utils} from "../../wailsjs/go/models";
import {BrowserOpenURL, LogInfo} from "../../wailsjs/runtime";
import ZiboInstallation = utils.ZiboInstallation;
import InstalledLivery = installer.InstalledLivery;
import AvailableLivery = installer.AvailableLivery;


function Zibo() {
    // const screens = useBreakpoint();
    const [running, setRunning] = useState(false);
    const [ziboDetails, setZiboDetails] = useState({} as ZiboInstallation);
    const [progressDetails, setProgressDetails] = useState("")
    const [installedLiveries, setInstalledLiveries] = useState([] as InstalledLivery[]);
    const [availableLiveries, setAvailableLiveries] = useState([] as AvailableLivery[]);
    useEffect(() => {
        const fetchDetails = async () => {
            const details = await FindZiboInstallationDetails();
            const liveries = await GetLiveries(details);
            setInstalledLiveries(liveries);
            setZiboDetails(details)
        };
        // Call it once immediately
        fetchDetails();
        // Set an interval to call it every 30 seconds
        const interval = setInterval(fetchDetails, 30000);

        // Clear the interval when the component is unmounted
        return () => clearInterval(interval);

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
        await BackupZiboInstallation(ziboDetails);
        const details = await FindZiboInstallationDetails();
        setZiboDetails(details)
        setRunning(false);
    }

    const handleRestore = async () => {
        setRunning(true);
        setProgressDetails("Restoring ...")
        await RestoreZiboInstallation(ziboDetails, "");
        const details = await FindZiboInstallationDetails();
        setZiboDetails(details)
        setRunning(false);
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
        await InstallZibo(ziboDetails, downloadInfo.path);
        const details = await FindZiboInstallationDetails();
        setZiboDetails(details)
        setRunning(false);
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
        await UpdateZibo(ziboDetails, downloadInfo.path);
        const details = await FindZiboInstallationDetails();
        setZiboDetails(details)
        setRunning(false);
    }

    return (
        <Card style={{
            minHeight: "100%",
        }}
        >
            <Spin spinning={running} tip={progressDetails}>
                <Row style={{minHeight: "100%"}}>
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
                            danger={ziboDetails.version !== ""}
                            style={{marginRight: "12px"}}
                            onClick={handleInstall}>{ziboDetails.version === "" ? "Install" : "Reinstall"}</Button>
                        <Button
                            type={"primary"}
                            style={{marginRight: "12px"}}
                            disabled={
                                ziboDetails.version === ziboDetails.remoteVersion ||
                                ziboDetails.version === ""
                            }
                            onClick={handleUpdate}>Update</Button>
                        <Button
                            type={"primary"}
                            style={{marginRight: "12px"}}
                            disabled={ziboDetails.version === ""}
                            onClick={handleBackup}>Backup</Button>
                        <Button
                            type={"primary"}
                            danger={true}
                            style={{marginRight: "12px"}}
                            disabled={ziboDetails.backupVersion === "N/A"}
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
                                    children: `${ziboDetails.version}`,

                                },
                                {
                                    key: '2',
                                    label: 'Current Version',
                                    children: `${ziboDetails.remoteVersion}`,
                                },
                                {
                                    key: '3',
                                    label: 'Backup Version',
                                    children: `${ziboDetails.backupVersion}`,
                                },
                                {
                                    key: '10',
                                    label: 'Installed Path',
                                    children: `${ziboDetails.path}`,
                                },

                            ]}
                        />
                    </Col>
                    <Divider/>
                    <Col span={24}>
                        <Tabs items={[
                            {
                                label: "Installed Liveries",
                                key: "installedLiveries",
                                children: <Table
                                    // title={() => "Liveries"}
                                    dataSource={installedLiveries}
                                    style={{overflow: "scroll"}}
                                    columns={[
                                        {
                                            title: 'Icon',
                                            dataIndex: 'icon',
                                            key: 'icon',
                                            render: (icon: string) => <img src={`data:image/png;base64,${icon}`}
                                                                           alt={"icon"}
                                                                           width={160} height={90}/>
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
                                    style={{overflow: "scroll"}}
                                    pagination={{pageSize: 25}}
                                    columns={[
                                        {
                                            title: 'Icon',
                                            dataIndex: 'icon',
                                            key: 'icon',
                                            render: (icon: string) => <img src={`data:image/png;base64,${icon}`}
                                                                           alt={"icon"}
                                                                           width={320} height={200}/>
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
