import React from 'react'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'react-chartjs-2';
import {useEffect, useState} from 'react';
import './App.css';
import {Cpu, Disk } from "../wailsjs/go/main/App";
ChartJS.register(ArcElement, Tooltip, Legend);


function App() {
    const [data, setData] = useState(0);
    const [data1, setData1] = useState([0,0]);

    const updateResultText = (newValue) => setData(newValue);

    const updateResultText1 = (newValue) => {
        setData1([newValue[0]/(1024*1024*1024), newValue[1]/(1024*1024*1024)]);
    };

    useEffect(() => {
        setInterval(() => {
            Cpu().then(updateResultText);
            Disk().then(updateResultText1);
        }, 2000);
    }, []);

    return (
        <div id="App">
            <>
            <div className="card mb-3 Card container">
                <div className="row g-0">
                    <div className="col-md-4">
                    </div>
                    <div className="col-md-8">
                        <div className="card-body">
                            <h2 className="card-title CardTexts">Proyecto 1</h2>
                            <h2 className="card-text CardTexts">Diego Andrés Obín Rosales</h2>
                            <h2 className="card-text CardTexts">201903865</h2>
                            <h2 className="card-text CardTexts"><small className="text-muted CardTexts">Sistemas operativos 2</small></h2>
                        </div>
                    </div>
                </div>
            </div>
            <div className='container ButtonsHome'>
                <h1>CPU</h1>
                <div>
                    <h3>{data.toFixed(2)}%</h3>
                </div>
                <h1>Memoria</h1>
                <div style={{display:"flex", justifyContent:"center"}}>
                    <table style={{border: "2px solid white"}}>
                        <thead>
                            <tr>
                            <th>Libre</th>
                            <th>En Uso</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                            <td>{data1[0].toFixed(2)}Gb</td>
                            <td>{data1[1].toFixed(2)}Gb</td>
                            </tr>
                        </tbody>
                        </table>
                </div>
            </div>
        </>
        </div>
    )
}

export default App
