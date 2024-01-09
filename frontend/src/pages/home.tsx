import React, {useEffect} from 'react';
import type {MenuProps} from "antd";
import {Tabs} from 'antd'
import {GetConfig} from "../../wailsjs/go/main/App";
import Zibo from "../components/zibo";
import Config from "../components/config";
import Backup from "../components/backup";


function Home() {
    // let allPlugins: utils.PluginConfig[] = [];

    const items: MenuProps['items'] = [
        {
            key: '1',
            label: 'TODO',
        },

    ];


    useEffect(() => {
        (async () => {
            const config = await GetConfig();
        })();

    }, []);


    return (
        <Tabs
            tabPosition={"left"}
            size={"large"}
            type={"card"}
            defaultActiveKey={"home"}
            items={
                [
                    {
                        label: "Zibo 737",
                        key: "home",
                        children: <Zibo/>,
                    },
                    {
                        label: "Configuration",
                        key: "configuration",
                        children: <Config/>,
                    },
                    {
                        label: "Backups",
                        key: "backups",
                        children: <Backup/>,
                    },
                ]
            }
            style={
                {
                    minHeight: "600px",
                    overflow: "scroll"
                }
            }
        />
    )
}

export default Home
