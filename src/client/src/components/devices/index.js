import React, { Component } from 'react';
import FirmataDevice from './FirmataDevice';
import GenericPlugin from './GenericPlugin';

export default class Devices extends Component {

  constructor(props) {
    super(props);
  }

  render() {
    let { device, plugin } = this.props;

    if (device && device.type === 'firmata') return <FirmataDevice device={device} />
    
    if (plugin && plugin.name) return <GenericPlugin plugin={plugin} />

    return null;
  }

}