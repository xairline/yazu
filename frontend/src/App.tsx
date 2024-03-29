import React, {AriaAttributes, DOMAttributes, useEffect, useState} from 'react';
import 'antd/dist/reset.css';
import './App.css';
import {Button, Col, Image, Layout, Modal, Row, Typography} from 'antd'
import {CheckXPlanePath, GetConfig, GetLatestVersion, GetVersion, IsXPlanePathConfigured} from "../wailsjs/go/main/App";
import {Link, Route, Routes} from "react-router-dom";
import useBreakpoint from "antd/es/grid/hooks/useBreakpoint";
import logo from './assets/images/logo-universal.png';
import Home from "./pages/home";
import Config from "./components/config";
import Paragraph from 'antd/es/typography/Paragraph';
import {BrowserOpenURL} from "../wailsjs/runtime";

const {Content,} = Layout;
const {Title} = Typography;

function App() {
    const screens = useBreakpoint();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isPathValid, setPathValid] = useState(false);
    const [xplanePath, setXplanePath] = useState("");
    const [localVersion, setLocalVersion] = useState("");
    const [remoteVersion, setRemoteVersion] = useState("");
    const showModal = () => {
        setIsModalOpen(true);
    };

    const handleOk = () => {
        IsXPlanePathConfigured().then((isPathValid) => {
            if (isPathValid) {
                setPathValid(isPathValid);
                setIsModalOpen(false);
            }
        });

    };

    useEffect(() => {
        (async () => {
            const config = await GetConfig();
            const isPathValid = await CheckXPlanePath(config.XPlanePath, []);
            const localVersion = await GetVersion();
            const remoteVersion = await GetLatestVersion();
            setLocalVersion(localVersion);
            setRemoteVersion(remoteVersion);
            setPathValid(isPathValid);
            if (!isPathValid) {
                console.log(JSON.stringify(config));
                showModal();
            } else {
                setXplanePath(config.XPlanePath);
            }
        })();

    }, []);


    return (
        <Layout className="layout app">
            <Row style={{background: "#006363"}}>
                <Col span={2}>
                    <Row style={{
                        display: "flex",
                        height: "100%",
                    }}>
                        <Link to={"/"}>
                            <Image src={logo}
                                   style={
                                       {
                                           maxHeight: "8vh",
                                           objectFit: "contain",
                                           margin: "12px 24px 12px"
                                       }
                                   }
                                   preview={false}
                            >

                            </Image>
                        </Link>
                    </Row>
                </Col>
                <Col span={6} offset={16}>
                    <Row style={{
                        display: "flex",
                        height: "100%",
                    }}>
                        <Paragraph
                            style={{
                                marginTop: "36px",
                                marginRight: "8px",
                                width: "100%", color: "white"
                            }}

                        >{localVersion}<Button
                            className="site-badge-count-109"
                            type={"primary"}
                            disabled={localVersion === remoteVersion}
                            style={{
                                backgroundColor: localVersion === remoteVersion ? '#52c41a' : '#1890ff',
                                marginLeft: "8px",

                            }}
                            onClick={() => {
                                BrowserOpenURL("https://github.com/xairline/yazu/releases/latest")
                            }}
                        >
                            {localVersion !== remoteVersion ? "Update Available" : "Up to date"}
                        </Button></Paragraph>

                    </Row>
                </Col>
            </Row>

            <Layout>
                <Content
                    style={{
                        padding: 24,
                        minHeight: "100%",
                        background: "white",
                        overflow: "hidden",
                    }}
                >
                    <Modal title="Configure X Plane Path"
                           open={isModalOpen}
                           onOk={handleOk}
                        // okButtonProps={{disabled: !isPathValid}}
                           cancelButtonProps={{hidden: true}}
                    >
                        <Config/>
                    </Modal>
                    <Routes>
                        <Route path={"/"} element={isPathValid && <Home/>}/>
                        {/*<Route path={"/dashboard"} element={<Dashboard/>}/>*/}
                        {/*<Route path="/callback" element={<Callback/>}/>*/}
                        {/*<Route*/}
                        {/*    key={'flight-logs'}*/}
                        {/*    path="/flight-logs/:id"*/}
                        {/*    element={<FlightLog/>}*/}
                        {/*/>*/}
                    </Routes>
                </Content>
            </Layout>
        </Layout>
    )
}

export default App

declare module 'react' {

    interface HTMLAttributes<T> extends AriaAttributes, DOMAttributes<T> {
        // extends React's HTMLAttributes
        directory?: string;        // remember to make these attributes optional....
        webkitdirectory?: string;
    }

}