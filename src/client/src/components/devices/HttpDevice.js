import React, { Component } from "react";
import Http from "../../icons/Http";
import ButtonRemote from "../ButtonRemote";
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
            <p>{device.description}</p>
          </div>
          <Http className="card-icon" />
        </div>
        <div className="actions">
          {device.actions &&
            device.actions.map((a, i) => (
              <ButtonRemote url={"/action"} data={{ name: a.name }} key={i}>
                {a.action.action}
              </ButtonRemote>
            ))}
        </div>
      </div>
    );
  }
}
