import React, { Component } from "react";
import ApiClient from "../tools/ApiClient";

export default class ButtonRemote extends Component {
  constructor(props) {
    super(props);
  }

  onPress = () => {
    if (this.props.data) {
      ApiClient.post(this.props.url, this.props.data)
        .then(res => {
          console.log("response click", res);
        })
        .catch(error => {
          console.log("error", error);
        });
    } else {
      ApiClient.get(this.props.url)
        .then(res => {
          console.log("response click", res);
        })
        .catch(error => {
          console.log("error", error);
        });
    }
  };

  render() {
    return <button onClick={this.onPress}>{this.props.children}</button>;
  }
}
