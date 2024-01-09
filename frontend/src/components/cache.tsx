import React, {useEffect, useState} from 'react';
import {GetBackups, GetCachedFiles} from "../../wailsjs/go/main/App";
import {Card} from "antd";


function Cache() {
    const [backups, setBackups] = useState({});
    const [cachedFiles, setCachedFiles] = useState({});
    useEffect(() => {
        (async () => {
            const backups = await GetBackups();
            const cachedFiles = await GetCachedFiles();
            setBackups(backups)
            setCachedFiles(cachedFiles)
        })();

    }, []);


    return (

        <>
            <Card style={{
                minHeight: "100%",
            }}
            >
                {JSON.stringify(backups)}
            </Card>
            <Card style={{
                minHeight: "100%",
            }}
            >
                {JSON.stringify(cachedFiles)}
            </Card>
        </>
    )
}

export default Cache
