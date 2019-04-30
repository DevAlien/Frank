import React, { Component } from "react";
import ApiClient from "../tools/ApiClient";
import Commands from "../components/commands";
import Gx from "gx";
export default class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      commands: []
    };
  }

  componentDidMount() {
    ApiClient.get("/commands")
      .then(res => {
        console.log("commands", res);
        this.setState({ commands: res });
      })
      .catch(error => {
        console.log("error", error);
      });
  }

  render() {
    let { commands } = this.state;
    return (
      <div>
        <h2>Home</h2>
        <div>
          <Gx col={4}>
            {commands.length > 0 && (
              <div>
                <h3>Commands</h3>
                {commands.map(d => <Commands command={d} />)}
              </div>
            )}
          </Gx>
        </div>
      </div>
    );
  }
}
