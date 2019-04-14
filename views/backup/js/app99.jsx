var { creatBrowserHistory, BrowserRouter, Route,Redirect,Link  } = ReactRouterDOM;
class App extends React.Component {

  setState() {
    try {
    var  data = JSON.parse(Cookies.get("Gate_user"));
    }
    catch(error){
     var  data = {"loggedIn": false};
    }
    // console.log(data);
    if (data.loggedIn) {
      this.logged= true;
    } else {
      this.logged= false;
  }
}

  componentWillMount() {
    this.setState();
  }
  render() {

    if (this.logged) {
      return (<LoggedIn />);
    } else {
      return ( <Home  LoggedIn={this.logged} /> );
    }
  }
}


class Home extends React.Component {
  render() {
    return (
      <div class="">
        <Navigation LoggedIn={this.props.LoggedIn} />
      </div>
    )
  }
}






class Navigation extends React.Component {
  render() {
    return (
      <nav class="navbar navbar-inverse">
        <div class="container-fluid">
          <div class="navbar-header">
            <a class="navbar-brand"  href="/">Modbus to MQTT Gateway Web Interface</a>
          </div>
          {this.props.LoggedIn ? <NavLoggedIn /> : <NavLogin /> }
          
        </div>
      </nav>
    )
  }
}



class NavLoggedIn extends React.Component {
  
  render() {
    return(
      <ul class="nav navbar-nav  navbar-right">
        <li class="dropdown"><a class="dropdown-toggle" data-toggle="dropdown" href="#">Hi  <span class="caret"></span></a>
          <ul class="dropdown-menu navbar-right">
            <li><Link to="/logout">Logout</Link></li>
            <li><Link to="/settings/view">Mqtt & Modbus Settings</Link></li>
            <li><Link to="/Dev/view">Devices</Link></li>
            <li><Link to="/modRead/view">Modbus Read Registers FC 0x03 </Link></li>
          </ul>
        </li>
      </ul>
    )
  }
}

class NavLogin extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      error: '',
      showPopup: false
    }
    this.commonChange = this.commonChange.bind(this);
    this.togglePopup = this.togglePopup.bind(this);
    this.getResponse = this.getResponse.bind(this);
    this.authenticate = this.authenticate.bind(this);
  }

  togglePopup() {
    this.setState({ showPopup: !this.state.showPopup });
  }

  async getResponse() {
    const headers = new Headers();
    headers.append('Content-type', 'application/json');

    const options = {
      method: 'POST',
      headers,
      body: JSON.stringify({
        UserName: this.state.username,
        Password: this.state.password
      })
    };

    const request = new Request("/login", options);
    const response = await fetch(request);
    console.log(await response.status);
    // const data = await response.json()
    //             .then(response => {  return response.error; })
    //             .catch(e => { return e; });
    const data = await response.json()
    return data;
  }
  authenticate() {
    this.getResponse().then((data) => {
      if (data.msg != "Done") {
        alert(data.msg);
        this.togglePopup()
        this.setState({ error: data.msg });
      } else if (data.msg === "Done") {
        // alert(data.msg);
        console.log("( <Redirect to");
        return ( <Redirect to="/"/> );
        
      }
    });
  }

  commonChange(event) {
    // console.log(this.state.username);
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  render() {
    return(
      <div  class="navbar-form navbar-right">
        <div class="form-group">
                <span style= {{color:'aliceblue'}} >Login ID</span>  <input type= "text" name="username"  onChange={this.commonChange} /> 
                <span style= {{color:'aliceblue'}} > Password </span> <input type= "password" name= "password" onChange={this.commonChange} /> 
        </div>
        
        <input type="submit" class="btn btn-success"  value="Login" onClick={this.authenticate} />
      </div >
    )
  }
}


ReactDOM.render(
  <BrowserRouter>
    <div>
      <Route path="/" component={App}></Route>
      <Route path="/App" component={App}></Route>
    </div>
  </BrowserRouter>,
  document.getElementById('app')
);