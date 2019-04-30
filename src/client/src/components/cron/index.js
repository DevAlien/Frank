import React, { Component } from "react";
import cronstrue from "cronstrue/i18n";
import ButtonRemote from "../ButtonRemote";
export default class Cron extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    this.getDescription(this.props.cron.cron_expression);
  }
  getDescription = command => {
    return cronstrue.toString(command, { locale: "it" });
  };

  render() {
    let { cron } = this.props;

    return (
      <div className="card">
        <div className="rich-area">
          <div className="content">
            <h2>{cron.description}</h2>
            <h4 style={{ margin: 5 }}>{cron.cron_expression}</h4>
            <h5 style={{ marginTop: 5 }}>
              {this.getDescription(cron.cron_expression)}
            </h5>
          </div>
        </div>
      </div>
    );
  }
}
