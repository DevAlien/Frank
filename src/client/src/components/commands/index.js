import React, { Component } from "react";

import ButtonRemote from "../ButtonRemote";
export default class Command extends Component {
  constructor(props) {
    super(props);
  }

  buildActions = (action, command, i) => {
    let found = findBrackets(command);

    if (found !== false) {
      if (action.matchingAction && action.matchingAction[found]) {
        let toReturn = [];
        for (var key in action.matchingAction[found]) {
          if (action.matchingAction[found].hasOwnProperty(key)) {
            toReturn.push(
              <ButtonRemote
                url={"/command?text=" + command.replace(/\{(.*?)\}/gm, key)}
                key={i + key}
              >
                {key}
              </ButtonRemote>
            );
          }
        }
        return toReturn;
      } else {
        return (
          <div>
            {command} <input />
          </div>
        );
      }
    } else {
      return (
        <ButtonRemote url={"/command?text=" + command} key={i}>
          {action.action}
        </ButtonRemote>
      );
    }
  };

  render() {
    let { command } = this.props;

    return (
      <div className="card">
        <div className="rich-area">
          <div className="content">
            <h2>{command.name}</h2>
          </div>
        </div>
        <div className="actions">
          {command.actions &&
            command.actions.map((a, i) => {
              return this.buildActions(a, command.commands[0], i);
            })}
        </div>
      </div>
    );
  }
}

const findBrackets = str => {
  console.log("dio1", str);
  const regex = /\{(.*?)\}/gm;
  let m,
    found = false;
  while ((m = regex.exec(str)) !== null) {
    if (m.index === regex.lastIndex) {
      regex.lastIndex++;
    }
    found = m[1];
  }

  return found;
};
