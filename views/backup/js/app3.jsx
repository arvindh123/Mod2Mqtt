var insert_point = document.querySelector('#container');
var { BrowserRouter, Route } = ReactRouterDOM;



class App extends React.Component {
  render() {
    return(
      <div>
        <h2>Route A</h2>
       
      </div>
    );
  }
}

ReactDOM.render(
  
  <BrowserRouter>
    <Route path="/" component={App}></Route>
    <Route path="/" component={App}></Route>
  </BrowserRouter>,
  document.getElementById('app')
);