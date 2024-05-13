import React, { useState } from "react";
import runCrawler from "./crawl_runner";

function CrawlerRunner() {
    const [credentials, setCredentials] = useState([]);
    const [formData, setFormData] = useState({
        ip: '',
        password: '',
        path_to_crawler: '',
        isHost: false,
        path_to_task_manager: ''
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData({
            ...formData,
            [name]: type === 'checkbox' ? checked : value
        });
    };

    const handleSubmit = () => {
        setCredentials(credentials.concat([formData]));
        setFormData({
            ip: '',
            password: '',
            path_to_crawler: '',
            isHost: false,
            path_to_task_manager: ''
        });
    };

    const handleDelete = (index) => {
        const updatedCredentials = credentials.filter((_, i) => i !== index);
        setCredentials(updatedCredentials);
    };

    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            width: '100vw',
            backgroundColor: 'gray',
            color: 'black',
        
        }}>
            <div style={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                height: '30%',
                width: '50%',
                margin: 'auto',
                backgroundColor: 'lightgray',
                color: 'black',
            }}>
                <input type="text" name="ip" placeholder="IP Address" value={formData.ip} onChange={handleChange} />
                <input type="password" name="password" placeholder="Password" value={formData.password} onChange={handleChange} />
                <input type="text" name="path_to_crawler" placeholder="Path to Crawler" value={formData.path_to_crawler} onChange={handleChange} />
                <label>
                    <input type="checkbox" name="isHost" checked={formData.isHost} onChange={handleChange} />
                    Is Host
                </label>
                <input type="text" name="path_to_task_manager" placeholder="Path to Task Manager" value={formData.path_to_task_manager} onChange={handleChange} />
                <button onClick={handleSubmit}>Add Credentials</button>
            </div>

            <div style={{
                display: 'flex',
                flexDirection: 'row',
                
                justifyContent: 'center',
                alignItems: 'center',
                height: '30vh',
                width: '50vw',
                margin: 'auto',
                backgroundColor: 'lightgray',
                color: 'black',
            }}>
                {credentials.map((credential, index) => (
                    <div key={index} style={{
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'center',
                        alignItems: 'center',
                        height: '30vh',
                        width: '30%',
                        margin: 'auto',
                        backgroundColor: 'lightgray',
                        color: 'black',
                    }}>
                        <ul>
                            <li>IP: {credential.ip}</li>
                            <li>Path to Crawler: {credential.path_to_crawler}</li>
                            <li>Is Host: {credential.isHost.toString()}</li>
                            <li>Path to Task Manager: {credential.path_to_task_manager}</li>
                        </ul>
                        <button onClick={() => handleDelete(index)}>Delete</button>

                    </div>
                ))}
            </div>

            <div style={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: 'center',
                height: '30vh',
                width: '50vw',
                margin: 'auto',
                backgroundColor: 'lightgray',
                color: 'black',
            }}>
                <button onClick={() => {
                    credentials.forEach(credential => {
                        runCrawler(credential);
                    });
                }} style={{
                    width: '100%',
                    height: '100%',
                    margin: 'auto',
                    backgroundColor: 'lightblue',
                    color: 'black',
                
                }}>Run Crawler System</button>
            </div>
        </div>
    );
}

export default CrawlerRunner;
