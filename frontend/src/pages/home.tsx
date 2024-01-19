import React, {useEffect, useState} from 'react';
import {Row, Spin, Tabs} from 'antd'
import {FindZiboInstallationDetails, GetConfig} from "../../wailsjs/go/main/App";
import Config from "../components/config";
import Backup from "../components/backup";
import Zibo from "../components/zibo";
import {utils} from "../../wailsjs/go/models";
import {RocketOutlined} from "@ant-design/icons";


function Home() {
    // let allPlugins: utils.PluginConfig[] = [];
    const [ziboDetails, setZiboDetails] = useState({} as utils.ZiboInstallation[]);
    useEffect(() => {
        (async () => {
            const config = await GetConfig();
            const details = await FindZiboInstallationDetails();
            setZiboDetails(details)
        })();

    }, []);


    return (
        (ziboDetails.length > 0 ?
            <Tabs
                tabPosition={"left"}
                size={"large"}
                type={"card"}
                defaultActiveKey={"home"}
                items={
                    [

                        ...(ziboDetails.map((ziboDetail: utils.ZiboInstallation) => {
                            return {
                                label:
                                    <Row><RocketOutlined
                                        style={{marginRight: 12}}/>{ziboDetail.path.split("/Aircraft/")[1].split("/plugins/")[0]}
                                    </Row>,
                                key: ziboDetail.path,
                                children: <Zibo installationDetails={ziboDetail}/>,
                            }
                        })),
                        {
                            label: "Backups",
                            key: "backups",
                            children: <Backup/>,
                        },
                        {
                            label: "Configuration",
                            key: "configuration",
                            children: <Config/>,
                        },
                    ]
                }
                style={
                    {
                        minHeight: "600px",
                        overflow: "scroll"
                    }
                }
            /> : <Spin spinning={true} tip={"Loading..."}/>)
    )
}

export default Home
