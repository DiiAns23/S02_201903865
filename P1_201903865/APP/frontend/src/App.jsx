import React, { Component } from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import ReactApexChart from "react-apexcharts";
import { useEffect, useState } from "react";
import "./App.css";
import { Cpu, Disk, Ram } from "../wailsjs/go/main/App";

ChartJS.register(ArcElement, Tooltip, Legend);

const Modal = ({ onClose }) => {
    return (
        <div style={{
            position: 'fixed',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            backgroundColor: 'white',
            padding: '2rem',
            zIndex: 1000,
        }}>
            <h2>Ventana emergente</h2>
            <p>Contenido de la ventana emergente...</p>
            <button onClick={onClose}>Cerrar</button>
        </div>
        );
};

const user = "";
const password = "";

const login= () => {
    console.log("Usuario y password: ", user, password)
    if (user == "admin" && password == "admin"){
        return true;
    }
    else{
        return false;
    }
}
class App extends Component {
    
    constructor(props) {
        super(props);

        this.state = {
        options_cpu: {
            chart: {
            id: "bar-chart",
            },
            xaxis: {
            categories: ["CPU"],
            },
            series: [
            {
                name: "Ocupado",
                data: [0],
            },
            {
                name: "Libre",
                data: [0],
            },
            ],
        },
        options_ram: {
            chart: {
            id: "bar-chart",
            },
            xaxis: {
            categories: ["RAM"],
            },
            series: [
            {
                name: "Ocupado",
                data: [0],
            },
            {
                name: "Libre",
                data: [0],
            },
            ],
        },
        options_disk: {
            chart: {
            id: "bar-chart",
            },
            xaxis: {
            categories: ["DISK"],
            },
            series: [
            {
                name: "Ocupado",
                data: [0],
            },
            {
                name: "Libre",
                data: [0],
            },
            ],
        },
        };
    }

    componentDidMount() {
        // Actualizar los datos cada segundo
        setInterval(() => {
        Cpu().then((data) => {
            this.setState({
            options_cpu: {
                ...this.state.options_cpu,
                series: [
                {
                    ...this.state.options_cpu.series[0],
                    data: [
                    Math.round(data)+ "%", // Valor actualizado en la primera serie
                    ...this.state.options_cpu.series[0].data.slice(1), // Mantener los valores originales en la primera serie
                    ],
                },
                {
                    ...this.state.options_cpu.series[1],
                    data: [
                    100 - Math.round(data)+ "%", // Valor actualizado en la primera serie
                    ...this.state.options_cpu.series[0].data.slice(1), // Mantener los valores originales en la primera serie
                    ],
                },
                ], // Mantener las demás series sin cambios
            },
            });
        });

        Ram().then((data) => {
            this.setState({
            options_ram: {
                ...this.state.options_ram,
                series: [
                {
                    ...this.state.options_ram.series[0],
                    data: [
                        Math.round(data[1]/(1024*1024*1024))+" GB" , // Valor actualizado en la primera serie
                    ...this.state.options_ram.series[0].data.slice(1), // Mantener los valores originales en la primera serie
                    ],
                },
                {
                    ...this.state.options_ram.series[1],
                    data: [
                        Math.round(data[0]/(1024*1024*1024))+" GB", // Valor actualizado en la primera serie
                    ...this.state.options_ram.series[0].data.slice(1), // Mantener los valores originales en la primera serie
                    ],
                },
                ], // Mantener las demás series sin cambios
            },
            });
        });

        Disk().then((data) => {
            this.setState({
            options_disk: {
                ...this.state.options_disk,
                series: [
                {
                    ...this.state.options_disk.series[0],
                    data: [
                        Math.round(data[1]/(1024*1024*1024))+" GB" , // Valor actualizado en la primera serie
                    ...this.state.options_disk.series[0].data.slice(1), // Mantener los valores originales en la primera serie
                    ],
                },
                {
                    ...this.state.options_disk.series[1],
                    data: [
                        Math.round(data[0]/(1024*1024*1024))+" GB", // Valor actualizado en la primera serie
                    ...this.state.options_disk.series[0].data.slice(1), // Mantener los valores originales en la primera serie
                    ],
                },
                ], // Mantener las demás series sin cambios
            },
            });
        });
        }, 1000);
    }

        
    render() {
        return (
        <div>
            {this.logged ? (
                <div>
                    <h1>Monitor de recursos</h1>
                </div>
                ) : (
                <div>
                    <h1>Monitor de recursos</h1>
                    <br />
                    <form action="">
                        <h2 for="user">Usuario</h2>
                        <input id={user} type="text" />
                        <h2 for="password">  Contraseña</h2>
                        <input id={password} type="password" />
                        <br />
                        <button type="submit" onClick={login()}>Iniciar sesión</button>
                    </form>
                    <h1>CPU</h1>
                    <ReactApexChart
                        options={this.state.options_cpu}
                        series={this.state.options_cpu.series}
                        type="bar"
                        height={400}
                    />
                    <br />
                    <h1>RAM</h1>
                    <ReactApexChart
                        options={this.state.options_ram}
                        series={this.state.options_ram.series}
                        type="bar"
                        height={400}
                    />
                    <br />
                    <h1>DISK</h1>
                    <ReactApexChart
                        options={this.state.options_disk}
                        series={this.state.options_disk.series}
                        type="bar"
                        height={400}
                    />
                </div>
                )}
        </div>
    );
    }
}

export default App;
