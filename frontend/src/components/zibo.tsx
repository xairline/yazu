import React, {useEffect, useState} from 'react';
import {Button, Card, Col, Descriptions, Divider, Row, Spin, Table} from 'antd'
// import {FindZiboInstallationDetails} from "../../wailsjs/go/main/App";
import useBreakpoint from "antd/es/grid/hooks/useBreakpoint";
import {
    BackupZiboInstallation,
    DownloadZibo,
    FindZiboInstallationDetails,
    GetDownloadDetails,
    InstallZibo,
    RestoreZiboInstallation,
    UpdateZibo
} from "../../wailsjs/go/main/App";
import {utils} from "../../wailsjs/go/models";
import ZiboInstallation = utils.ZiboInstallation;


function Zibo() {
    // const screens = useBreakpoint();
    const [running, setRunning] = useState(false);
    const [ziboDetails, setZiboDetails] = useState({} as ZiboInstallation);
    const [progressDetails, setProgressDetails] = useState("")
    useEffect(() => {
        (async () => {
            const details = await FindZiboInstallationDetails();
            setZiboDetails(details)
        })();

    }, []);

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
        await RestoreZiboInstallation(ziboDetails);
        const details = await FindZiboInstallationDetails();
        setZiboDetails(details)
        setRunning(false);
    }
    const handleInstall = async () => {
        setRunning(true);
        const isDownloading = await DownloadZibo(true);
        while (isDownloading) {
            const downloadDetails = await GetDownloadDetails(false);
            setProgressDetails(`${downloadDetails.toFixed(2)}%`)
            if (downloadDetails === 100) {
                break;
            }
        }
        setProgressDetails("Installing ...")
        await InstallZibo(ziboDetails);
        const details = await FindZiboInstallationDetails();
        setZiboDetails(details)
        setRunning(false);
    }

    const handleUpdate = async () => {
        setRunning(true);
        const isDownloading = await DownloadZibo(true);
        while (isDownloading) {
            const downloadDetails = await GetDownloadDetails(true);
            setProgressDetails(`${downloadDetails.toFixed(2)}%`)
            if (downloadDetails === 100) {
                break;
            }
        }
        setProgressDetails("Updating ...")
        await UpdateZibo(ziboDetails);
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
                        <Table
                            title={() => "Liveries"}
                        />
                    </Col>
                </Row>
            </Spin>
        </Card>
    )
}

export default Zibo
