import React, { useState } from "react";
import axios from "axios";
import "./styles.css"
function CrawlerRunner() {
    const [credentials, setCredentials] = useState([]);
    const [formData, setFormData] = useState({
        ip: '',
        username: '',
        password: '',
        path_to_crawler: '',
        is_host: false,
        path_to_task_manager: '',
    });
    const [esCredentials, setEsCredentials] = useState({
        cloud_id: '',
        api_key: '',
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        if (name === "cloud_id" || name === "api_key") {
            setEsCredentials({
                ...esCredentials,
                [name]: value,
            });
        } else {
            setFormData({
                ...formData,
                [name]: type === 'checkbox' ? checked : value
            });
        }
    };

    const handleSubmit = () => {
        setCredentials(credentials.concat([formData]));
        setFormData({
            ip: '',
            username: '',
            password: '',
            path_to_crawler: '',
            path_to_task_manager: '',
            is_host: false,
        });
    };

    const handleEsSubmit = async () => {
        await axios.post("http://localhost:3000/api/api_credentials", esCredentials);
        setEsCredentials({
            cloud_id: '',
            api_key: '',
        });
        alert("Elasticsearch credentials updated!");
    };

    const handleDelete = (index) => {
        const updatedCredentials = credentials.filter((_, i) => i !== index);
        setCredentials(updatedCredentials);
    };

    const runCrawlers = async () => {
        let host = credentials.find(credential => credential.is_host);
        let workers = credentials.filter(credential => !credential.is_host);
        
        let data = {
            ip: host.ip,
            username: host.username,
            password: host.password,
            path_to_crawler: host.path_to_crawler,
            path_to_task_manager: host.path_to_task_manager,
            is_host: true,
            id: 1,
            host_ip: host.ip,
            worker_num: workers.length,
        }
        await axios.post("http://localhost:3000/api/execute_ssh", data);

        for (let i = 2; i < workers.length + 2; i++) {
            data = {
                ip: workers[i - 2].ip,
                username: workers[i - 2].username,
                password: workers[i - 2].password,
                path_to_crawler: workers[i - 2].path_to_crawler,
                path_to_task_manager: workers[i - 2].path_to_task_manager,
                is_host: false,
                id: i,
                host_ip: host.ip,
                worker_num: workers.length,
            };
            await axios.post("http://localhost:3000/api/execute_ssh", data);
        }
    };

    return (
        <div className="main">
            {/* Elasticsearch Credentials Form */}
            <div className="third">
                <input type="text" name="cloud_id" placeholder="Elasticsearch Cloud ID" value={esCredentials.cloud_id} onChange={handleChange} className="text-input"/>
                <input type="text" name="api_key" placeholder="Elasticsearch API Key" value={esCredentials.api_key} onChange={handleChange} className="text-input"/>
                <button onClick={handleEsSubmit} class = "btn-elastic">Set Elasticsearch Credentials</button>
            </div>
            <div className="first">
                <input type="text" name="ip" placeholder="IP Address" value={formData.ip} onChange={handleChange} className="text-input"/>
                <input type="text" name="username" placeholder="Username" value={formData.username} onChange={handleChange}className="text-input" />
                <input type="password" name="password" placeholder="Password" value={formData.password} onChange={handleChange} className="pass-input"/>
                <input type="text" name="path_to_crawler" placeholder="Path to Crawler" value={formData.path_to_crawler} onChange={handleChange} className="text-input"/>
                <input type="text" name="path_to_task_manager" placeholder="Path to Task Manager" value={formData.path_to_task_manager} onChange={handleChange} className="text-input"/>
                <div className="bottom">
                    <input type="checkbox" name="is_host" checked={formData.is_host} onChange={handleChange} class = "checkbox"/>
                    <p className="subtitle">Is Host</p>
                    <button onClick={handleSubmit} className = "btn">Add Credentials</button>
                </div>
            </div>
            <div className="second">
                {credentials.map((credential, index) => (
                    <div key={index} className="list-objs">
                        <ul>
                            <li>IP: {credential.ip}</li>
                            <li>Username: {credential.username}</li>
                            <li>Path to Crawler: {credential.path_to_crawler}</li>
                            <li>Is Host: {credential.is_host.toString()}</li>
                            <li>Path to Task Manager: {credential.path_to_task_manager}</li>
                        </ul>
                        <button onClick={() => handleDelete(index)} class = "btn-delete">Delete</button>

                    </div>
                ))}
            </div>

            <div className="input-cont">
                <button onClick={() => 
                    runCrawlers(credentials)
                } className="btn-run">Run Crawler System</button>
            </div>
        </div>
    );
}

export default CrawlerRunner;
