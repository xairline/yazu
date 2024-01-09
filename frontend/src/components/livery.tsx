import React, {useEffect, useState} from 'react';
import {Card, Col, Input, Row} from 'antd'
import {CheckXPlanePath, OpenDirDialog} from "../../wailsjs/go/main/App";
import useBreakpoint from "antd/es/grid/hooks/useBreakpoint";
import {CheckCircleTwoTone, CloseCircleTwoTone} from "@ant-design/icons";

interface Props {
    aircraft: string
}

const Livery: React.FC<Props> = (props: Props) => {
    const screens = useBreakpoint();
    const [isPathValid, setPathValid] = useState(false);
    const [xplanePath, setXplanePath] = useState("");
    useEffect(() => {
        (async () => {
            const isPathValid = await CheckXPlanePath(xplanePath);
            setPathValid(isPathValid);
        })();

    }, [xplanePath]);

    const handleFolderInputClick = async () => {
        const path = await OpenDirDialog()
        setXplanePath(path);

    };

    return (
        <Card>
            <Row style={{
                display: "flex",
                // justifyContent: !screens.sm ? "center" : "flex-end",
                alignItems: "center",
                height: "100%",
                marginTop: "10%"
                // marginRight: !screens.sm ? "12px" : "24px",
                // marginLeft: !screens.sm ? "12px" : "24px",
            }}>
                <Col flex={"auto"}>X Plane Path:</Col>
                <Col flex={"auto"}>
                    <Input value={xplanePath}
                           onClick={handleFolderInputClick}>

                    </Input>
                </Col>
                <Col flex={"auto"} style={{marginLeft: "24px"}}>
                    {
                        isPathValid ?
                            <CheckCircleTwoTone twoToneColor="#52c41a"/> :
                            <CloseCircleTwoTone twoToneColor="#eb2f96"/>
                    }
                </Col>
            </Row>
        </Card>
    )
}

export default Livery
