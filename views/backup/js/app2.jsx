var Router = ReactRouterDOM.Router;
var Route = ReactRouterDOM.Route;
var IndexRoute = ReactRouterDOM.IndexRoute;
var Link = ReactRouterDOM.Link;
var browserHistory = ReactRouterDOM.createBrowserHistory;

class RouteA extends React.Component {
  btnClickB() {
    browserHistory.push('/b');
  }
  btnClickC() {
    browserHistory.push('/c');
  }

  render() {
    return(
      <div>
        <h2>Route A</h2>
        <div>
          <button className="btn btn-default" onClick={ this.btnClickB }>Goto B</button>
          <button className="btn btn-default" onClick={ this.btnClickC }>Goto C</button>
        </div>
      </div>
    );
  }
}

class RouteB extends React.Component {
  btnClickA() {
    browserHistory.push('/a');
  }

  btnClickC() {
    browserHistory.push('/c');
  }

  render() {
    return(
      <div>
        <h2>Route B</h2>
        <div>
          <button className="btn btn-default" onClick={ this.btnClickA }>Goto A</button>
          <button className="btn btn-default" onClick={ this.btnClickC }>Goto C</button>
        </div>
      </div>
    );
  }
}

class RouteC extends React.Component {
  btnClickA() {
    browserHistory.push('/a');
  }

  btnClickB() {
    browserHistory.push('/b');
  }

  render() {
    return(
      <div>
        <h2>Route C</h2>
        <div>
          <button className="btn btn-default" onClick={ this.btnClickA }>Goto A</button>
          <button className="btn btn-default" onClick={ this.btnClickB }>Goto B</button>
        </div>
      </div>
    );
  }
}

ReactDOM.render(
  <Router history={browserHistory}>
    <Route path="/" component={RouteA} />
    <Route path="/a" component={RouteA} />
    <Route path="/b" component={RouteB} />
    <Route path="/c" component={RouteC} />
    <Route path="*" component={RouteA} />
  </Router>,
  document.getElementById('app')
);