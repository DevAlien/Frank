import React, { Component } from 'react';
import ApiClient from '../tools/ApiClient';
import Devices from '../components/devices';

export default class Settings extends Component {

  constructor(props) {
    super(props);

    this.state = {
      devices: [],
      plugins: []
    }
  }

  componentDidMount() {
    ApiClient.get('/devices?data=full').then(res => {
      console.log('devices', res);
      this.setState({devices: res});
    }).catch( error => {
      console.log('error', error)
    })

    ApiClient.get('/plugins').then(res => {
      console.log('plugins', res);
      this.setState({plugins: res});
    }).catch( error => {
      console.log('error', error)
    })
  }

  render() {
    let { devices, plugins } = this.state;

    return (
      <div>
        <h2>Settings</h2>
        {devices.length > 0 && <div>
          <h3>Devices</h3>
          {devices.map(d => <Devices device={d} />)}
        </div>}
        {plugins.length > 0 && <div>
          <h3>Plugins</h3>
          {plugins.map(p => <Devices plugin={p} />)}
        </div>}
      </div>
    );
  }

}