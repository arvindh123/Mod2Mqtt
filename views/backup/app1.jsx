var { BrowserRouter, Route,Redirect } = ReactRouterDOM;
class App extends React.Component {
  
  setStatee() {
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
    this.setStatee();
  }
  render() {

    console.log(this.props.location);

    if (!this.state.isLoggedin) {

      return (  <div class="">
      <Navigation />
      <div class="jumbotron text-center">
        <h1 >Login</h1>
        <p>Sign in to get access </p>
        {this.state.showPopup ? <Popup heading="Login Failed" text={this.state.error} btnstyle="btn-danger" closePopup={this.togglePopup.bind(this)} btntext="Close" /> : null}
        <label>User Name:
                  <input type="text" name="username" onChange={this.commonChange} />
        </label>
        <label>Password:
                  <input type="password" name="password" onChange={this.commonChange} />
        </label><br />

        <input type="submit" value="Submit" onClick={this.authenticate} />

      </div>
     
    </div>
      )
    } else {
      return <LoggedIn/>;
    }
  }
}

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      error: '',
      showPopup: false,

      isLoggedin:false
    }
    this.authenticate = this.authenticate.bind(this);
    this.commonChange = this.commonChange.bind(this);
    this.getResponse = this.getResponse.bind(this);
    this.togglePopup = this.togglePopup.bind(this);
  }
  commonChange(event) {
    // console.log(this.state.username)
    this.setState({
      [event.target.name]: event.target.value
    });
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
  togglePopup() {
    this.setState({ showPopup: !this.state.showPopup });
  }

  authenticate() {
    this.getResponse().then((data) => {
      if (data.msg != "Done") {
        // alert(data.msg);
        this.togglePopup()
        this.setState({ error: data.msg });
      } else if (data.msg === "Done") {
        alert(data.msg);
        
        this.setState({ isLoggedin: true });
         
      }
    });
  }

}


class LoggedIn extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      jokes: []
    }
  }
  serverRequest() {
    $.get("http://localhost:8080/user/users", res => {
      this.setState({
        users: res
      });
    });
  }

  componentDidMount() {
    this.serverRequest();
  }
  render() {

    return (
      <div className="container">
        <div className="col-lg-12">
          <br />
          <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
          <h2>Users</h2>


          <div className="row">
            {/* {this.state.users.map(function (user, i) {
              return (<Joke key={i} user={user} />);
            })} */}
          </div>
        </div>
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
          <NavLogin />
          
        </div>
      </nav>
    )
  }
}

class NavLogin extends React.Component {
  render() {
    return(
      <ul class="nav navbar-nav  navbar-right">
        <li class="dropdown"><a class="dropdown-toggle" data-toggle="dropdown" href="#">Hi  <span class="caret"></span></a>
          <ul class="dropdown-menu navbar-right">
            <li><a href="/logout">Logout</a></li>
            <li><a href="/settings/view">Mqtt & Modbus Settings</a></li>
            <li><a href="/Dev/view">Devices</a></li>
            <li><a href="/modRead/view">Modbus Read Registers FC 0x03 </a></li>
          </ul>
        </li>
      </ul>
    )
  }
}

ReactDOM.render(
  <BrowserRouter>
    <div>
      <Route path="/" component={App}></Route>
      <Route path="/login" component={LoggedIn}></Route>
      
    </div>
  </BrowserRouter>,
  document.getElementById('app')
);