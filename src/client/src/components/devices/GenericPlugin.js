import React, { Component } from 'react';
import Radio from '../../icons/Radio';

export default class GenericPlugin extends Component {

  constructor(props) {
    super(props);
  }

  render() {
    let { plugin } = this.props;

    return (
      <div className="card">
        <div className="rich-area">
          <div className="content">
            <h2>{plugin.name}</h2>
            <p>{plugin.name} Descrizione</p>
          </div>
          <Radio className="card-icon"/>
        </div>
        <div className="actions">
          {plugin.actions && plugin.actions.map(a => <button>{a.action.action}</button>)}
        </div>
      </div>
    );
  }

}