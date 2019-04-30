import React, { Component } from "react";
import ApiClient from "../tools/ApiClient";
import Devices from "../components/devices";
import Cron from "../components/cron";
import Gx from "gx";
export default class Settings extends Component {
  constructor(props) {
    super(props);

    this.state = {
      devices: [],
      plugins: [],
      crons: []
    };
  }

  componentDidMount() {
    ApiClient.get("/devices?data=full")
      .then(res => {
        console.log("devices", res);
        this.setState({ devices: res });
      })
      .catch(error => {
        console.log("error", error);
      });

    ApiClient.get("/plugins")
      .then(res => {
        console.log("plugins", res);
        this.setState({ plugins: res });
      })
      .catch(error => {
        console.log("error", error);
      });

    ApiClient.get("/crons")
      .then(res => {
        console.log("crons", res);
        this.setState({ crons: res });
      })
      .catch(error => {
        console.log("error", error);
      });
  }

  render() {
    let { devices, plugins, crons } = this.state;

    return (
      <div>
        <h2>Settings</h2>
        <div>
          <Gx col={4}>
            {devices.length > 0 && (
              <div>
                <h3>Devices</h3>
                {devices.map(d => <Devices device={d} />)}
              </div>
            )}
            {/* {plugins.length > 0 && (
              <div>
                <h3>Plugins</h3>
                {plugins.map(p => <Devices plugin={p} />)}
              </div>
            )} */}
          </Gx>
          <Gx col={4}>
            {crons.length > 0 && (
              <div>
                <h3>Crons</h3>
                {crons.map(d => <Cron cron={d} />)}
              </div>
            )}
          </Gx>
        </div>
      </div>
    );
  }
}
