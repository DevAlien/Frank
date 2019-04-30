import React, { Component } from "react";
import Arduino from "../../icons/Arduino";
export default class FirmataDevice extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    let { device } = this.props;

    return (
      <div className="card">
        <div className="rich-area">
          <div className="content">
            <h2>{device.name}</h2>
            <p>{device.name} Descrizione</p>
          </div>
          <Arduino className="card-icon" />
        </div>
        <div className="actions">
          {device.actions &&
            device.actions.map((a, i) => (
              <button key={i}>{a.action.action}</button>
            ))}
        </div>
      </div>
    );
  }
}
