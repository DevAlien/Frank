import React, { Component } from "react";
import FirmataDevice from "./FirmataDevice";
import HttpDevice from "./HttpDevice";
import GenericPlugin from "./GenericPlugin";

export default class Devices extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    let { device, plugin } = this.props;

    if (device && device.type === "firmata")
      return <FirmataDevice device={device} />;
    if (device && device.type === "http") return <HttpDevice device={device} />;

    if (plugin && plugin.name) return <GenericPlugin plugin={plugin} />;

    return null;
  }
}
