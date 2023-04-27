<script setup>
import {GetCPUUsage, GetDISKUsage, GetRAMUsage, GetUSBPorts, DeshabilitarUSB, HabilitarUSB, UsbLogs, AbrirArchivo} from '../../wailsjs/go/sys/Stats'
</script>

<template>
  <main>
    <div v-if="mostrarUsb">
      <form>
        <br/>
        <button @click.prevent="logout">Cerrar sesion</button>
      </form>
      <form>
        <br/>
        <button @click.prevent="logs">Archivo de Bitacora</button>
      </form>
      <div>
        <h2>Puertos en uso</h2>
        <table class="table custom-table align-center">
          <thead>
            <tr>
              <th>Puerto</th>
              <th>Nombre del dispositivo</th>
              <th>Device</th>
              <th>VendorID</th>
              <th>ProductID</th>
              <th>Accion</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(usb, index) in usbports" :key="index">
              <td>{{ usb.port_name }}</td>
              <td>{{ usb.name }}</td>
              <td>{{ usb.device }}</td>
              <td>{{ usb.vid }}</td>
              <td>{{ usb.pid }}</td>
              <td><button @click.prevent="des_habilitarUSB(usb.vid, usb.pid, index)">{{ statusArray[index] ? 'Deshabilitar' : 'Habilitar' }}</button></td>
            </tr>
          </tbody>
        </table>
      </div>

    </div>
    <div id="chart"  class="row" v-if="mostrarGraficas">
      <div class="column">
        <br/>
        <div>Manejo de USB</div>
        <br/>
        <form>
          <div><label for="user">Usuario</label></div>
          <div><input id="user" v-model="user" required></div>
          <div><label for="password">Contraseña</label></div>
          <div><input type="password" id="password" v-model="password" required></div>
          <button type="submit" @click.prevent="login">Ingresar</button>
        </form>
      </div>
      <div class="column">
        <apexchart
          ref="cpu"
          width="350"
          type="bar"
          :options="CPU_chartOptions"
          :series="CPU_series"
        ></apexchart>
      </div>
      <div class="column">
        <apexchart
          ref="disk"
          width="350"
          type="bar"
          :options="DISK_chartOptions"
          :series="DISK_series"
        ></apexchart>
      </div>
      <div class="column">
        <apexchart
          ref="ram"
          width="350"
          type="bar"
          :options="RAM_chartOptions"
          :series="RAM_series"
        ></apexchart>
      </div>
    </div>
  </main>
</template>


<script>
export default {
  name: "CPU",
  data: function() {
    return {
      mostrarGraficas: true,
      mostrarUsb: false,
      CPU_chartOptions: {
        chart: {
          id: 'vuechart-example'
        },
        xaxis: {
          categories: ['CPU (%)']
        }
      },
      CPU_series: [
        {
          name: "Ocupado",
          data: [50],
        },
        {
          name: "Libre",
          data: [50],
        }
      ],  
      DISK_chartOptions: {
        chart: {
          id: 'disk'
        },
        xaxis: {
          categories: ['Disco (GB)']
        }
      },
      DISK_series: [
        {
          name: "Ocupado",
          data: [50],
        },
        {
          name: "Libre",
          data: [50],
        }
      ],
      RAM_chartOptions: {
        chart: {
          id: 'ram'
        },
        xaxis: {
          categories: ['RAM (GB)']
        }
      },
      RAM_series: [
      {
          name: "Ocupado",
          data: [50],
        },
        {
          name: "Libre",
          data: [50],
        }
      ],
      user: '',
      password: '',
      usbports: [],
      statusArray: [],
      usbports_free: [],
      statusArray_free: [],
    };
  },
  mounted(){
    GetCPUUsage().then(cpu_usage => {
      this.CPU_series = [cpu_usage.avg, 100 - cpu_usage.avg]
    })
    GetDISKUsage().then(disk_usage => {
      this.DISK_series = [disk_usage.used, disk_usage.free]
    })
    GetRAMUsage().then(ram_usage => {
      this.RAM_series = [ram_usage.used, ram_usage.free]
    })
    this.updateChart();
  },
  methods: {
    updateChart() {
      setInterval(() => {
        GetCPUUsage().then(cpu_usage => {
          this.CPU_series = [
            {
              name: "Ocupado",
              data: [cpu_usage.avg],
            },
            {
              name: "Libre",
              data: [100 - cpu_usage.avg],
            }
          ]
        })
      }, 1000);
      setInterval(() => {
        GetDISKUsage().then(disk_usage => {
          console.log('disk_usage', disk_usage)
          this.DISK_series =  [
            {
              name: "Ocupado",
              data: [Math.round(disk_usage.used/(1024*1024*1024))],
            },
            {
              name: "Libre",
              data: [Math.round(disk_usage.free/(1024*1024*1024))],
            }
          ]
        })
      }, 5000);
      setInterval(() => {
        GetRAMUsage().then(ram_usage => {
          console.log('ram_usage', ram_usage)
          this.RAM_series =  [
            {
              name: "Ocupado",
              data: [Math.round(ram_usage.used/(1024*1024*1024))],
            },
            {
              name: "Libre",
              data: [Math.round(ram_usage.free/(1024*1024*1024))],
            }
          ]
        })
      }, 1000);
    },
    login() {
      console.log(`Usuario: ${this.user}, Password: ${this.password}`)
      if (this.user == "obin"){
        alert("Bienvenido")
        this.mostrarUsb = true
        this.mostrarGraficas =  false
        UsbLogs();
        GetUSBPorts().then(usb_ports => {
          this.usbports = usb_ports
          this.statusArray = this.usbports.map(port => port.status);
        })
      }else{
        alert("Error, usuario es 'obin' y la contraseña es la del administrador de linux")
      }
    },
    logout() {
      this.mostrarUsb = false
      this.mostrarGraficas =  true
      this.user = ''
      this.password = ''
    },
    logs(){
      AbrirArchivo();
    },
    des_habilitarUSB(vid, pid, index){
      if (this.statusArray[index]){
        DeshabilitarUSB(vid, pid, this.password).then(resp => {
          alert(resp)
        }).catch(error =>{
          alert(error)
        })
      } else{
        HabilitarUSB(vid, pid, this.password).then(resp => {
          alert(resp)
        }).catch(error =>{
          alert(error)
        })
      }
      this.statusArray[index] = !this.statusArray[index];
    }
  }
};
</script>

<style scoped>
  .column {
    float: left;
    width: 50%;
  }

  /* Clear floats after the columns */
  .row:after {
    content: "";
    display: table;
    clear: both;
  }

  .custom-table {
    border-collapse: collapse;
  }

  .custom-table th,
  .custom-table td {
    border: 1px solid black;
    padding: 10px;
  }
</style>
