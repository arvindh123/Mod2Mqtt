var { BrowserRouter, Link, withRouter, Route, Redirect } = ReactRouterDOM;
var Router = BrowserRouter;
// var {createReactClass} = ReactRouterDOM;


class AuthExample extends React.Component {
  componentWillMount() {
    WsClient.WsConnect();
  }
  render() {

    return (
      <Router basename="/#/">
        <div>
          <AuthButton />
          <Redirect to="/home" />
          <Route path="/home" component={Home} />
          <Route path="/login" component={NavLogIn} />
          <PrivateRoute path="/usermanager" component={UserManager} />
          <PrivateRoute path="/Interface" component={Interface} />
          <PrivateRoute path="/modregs" component={ModbusRegsiters} />
          <PrivateRoute path="/mqtt" component={Mqtt} />
          <PrivateRoute path="/devmodels" component={DeviceModels} />
          <PrivateRoute path="/devices" component={Devices} />
          <PrivateRoute path="/AdditionalFeatures" component={AdditionalFeatures} />
         
        </div>
      </Router>
    );
  }
};
const PostGet = {
  async Get(path,metho) {
      const headers = new Headers();
      headers.append('Content-type', 'application/json');
      const options = {
            method: metho,
            headers,
      };
      const request = new Request(path, options);
      const response = fetch(request);
      // const data = await response.json();
      return response;
  },
  async Post(path,metho,bod) {
    const headers = new Headers();
    headers.append('Content-type', 'application/json');
    const options = {
          method: metho,
          headers,
          body : JSON.stringify(bod)
    };
    const request = new Request(path, options);
    const response = fetch(request);
    // const data = await response.json();
    // console.log(response);
    return response;
  }
}

const WsClient = {
  conn :'',
  WsConnect(){
    if (this.conn.hasOwnProperty('readyState')){
      if (this.conn.readState !==0) {
        this.conn = new  WebSocket('ws://' + window.location.host + '/ws');

      }
    }else{
      this.conn = new  WebSocket('ws://' + window.location.host + '/ws');
    }
    
    // conn.
    // return conn
  }, 
  WsClose(){
    this.conn.close()
  },
  WsSendCmd(cmd){

    this.conn.send( JSON.stringify({"Cmd"  : cmd})  )
  },
  WsSend(msg){
    this.conn.send( JSON.stringify(msg))
  }
}
const RAuth = {
  isAuthenticated: false,
  async getResponse(username,password) {
    const headers = new Headers();
    headers.append('Content-type', 'application/json');
    const options = {
          method: 'POST',
          headers,
          body: JSON.stringify({
            UserName: username,
            Password: password
          })
        }
    const request = new Request("/login", options);
    const response = await fetch(request);
    const data = await response.json()
    return data;
    },

  authenticate(cb,username,password) {
    this.getResponse(username,password).then((data) => {
      
      if (data.msg[0].content != "Done") {
        alert(data.msg[0].content);
        // console.log(data.msg[0]);
        // console.log(data.msg[0].content);
        // this.togglePopup()
        // this.setState({ error: data.msg });
        this.isAuthenticated = false;
      } else if (data.msg[0].content  === "Done") {
        // alert(data.msg);
        this.isAuthenticated = true;
        setTimeout(cb, 100); 
      }
    });
    // console.log("Error")
    // this.isAuthenticated = true;
    // setTimeout(cb, 100); // fake async
  },

  async getSingout() {
    const response  = await fetch("/logout");
    const data = await response.json()         ;     
    return data;
  },

  signout(cb) {
    
    this.getSingout().then(data => {
            if (data.msg[0].Content !== "Done") { 
              this.isAuthenticated = false;
              
            }else { 
              this.isAuthenticated = false;
            } 
          })
          .catch(e => { console.log(e); })
      // if  (!this.isAuthenticated) {
      //   setTimeout(cb, 100);
      // }

      // this.isAuthenticated = false;
      setTimeout(cb, 100);
  }
  
};

const AuthButton = withRouter(
  ({ history }) =>
    RAuth.isAuthenticated ? (
      // <p>
      <NavLoggedIn
        PassFunc={() => {
          RAuth.signout(() => history.push("/home"));
        }} />
    ) : (
        <div>
          <NavLogIn location={window.location} />
          {/* <h1> Your not logged in </h1> */}
        </div>
      )
);



class NavLoggedIn extends React.Component {
  render() {
    return (
      <nav class="navbar navbar-inverse">
        <div class="container-fluid">
          <div class="navbar-header">
            <Link class="navbar-brand" to="/home">Modbus to MQTT Gateway Web Interface</Link>
          </div>

          <ul class="nav navbar-nav">
          <li>
              <Link to="/devices">Devices</Link>
            </li>
            <li>
              <Link to="/devmodels">Device Models</Link>
            </li>
          <li>
              <Link to="/interface">Interface</Link>
            </li>
            <li>
              <Link to="/modregs">Modbus Regsiters</Link>
            </li>
            
            <li>
              <Link to="/mqtt">Mqtt</Link>
            </li>
            <li>
              <Link to="/AdditionalFeatures">Additional Features</Link>
            </li>

            
            
          </ul>

          <ul class="nav navbar-nav  navbar-right">
            <li class="dropdown"><a class="dropdown-toggle" data-toggle="dropdown" href="#">Hi User</a>
              <ul class="dropdown-menu navbar-right">
                <li><a href="#" onClick={this.props.PassFunc} > Logout</a></li>
                <li> <Link to="/usermanager">UserManager</Link> </li> 
              </ul>
            </li>
          </ul>
        </div>
        {/* <h3> Hello</h3> */}
      </nav>
    
    );
  }
}


class NavLogIn extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: ''
    }
    this.commonChange = this.commonChange.bind(this)
  }
  commonChange(event) {
    // console.log(this.state.username)
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  state = { redirectToReferrer: false };
  login = () => {
    RAuth.authenticate(() => {
      this.setState({ redirectToReferrer: true });
    }, this.state.username, this.state.password);
  };
  render() {
    let { redirectToReferrer } = this.state;
    // console.log();    
    // let { from } = this.props.location || { from: { href: "/" } };
    // console.log({ from } );    
    if (redirectToReferrer) return <Redirect to={this.props.location.href.split("#")[1]}  />;
    return (
      <nav class="navbar navbar-inverse">
        <div class="container-fluid">
          <div class="navbar-header">
            <Link class="navbar-brand" to="/home">Modbus to MQTT </Link>
          </div>

          <ul class="nav navbar-nav">
          <li>
              <Link to="/devices">Devices</Link>
            </li>
            <li>
              <Link to="/devmodels">Device Models</Link>
            </li>
          <li>
              <Link to="/interface">Interface</Link>
            </li>
            <li>
              <Link to="/modregs">Modbus Regsiters</Link>
            </li>
            
            <li>
              <Link to="/mqtt">Mqtt</Link>
            </li>
            <li>
              <Link to="/AdditionalFeatures">Additional Features</Link>
            </li>
          </ul>

          <div class="navbar-form navbar-right" >
            <div class="form-group">
              <span style={{ color: "aliceblue" }} >Login ID</span>
              <input type="text" name="username" onChange={this.commonChange}/>
              <span style={{ color: "aliceblue" }} > Password </span>
              <input type="password" name="password" onChange={this.commonChange}/>
            </div>
            {/* <input type="submit" onClick={this.login} class="btn btn-success" value = "Login" /> */}
            <button class="btn btn-success" onClick={this.login}>Log in</button>
            {/* <button class="btn btn-success" >Log in</button> */}
          </div>
        </div>
      </nav>
    );
  }
}


function PrivateRoute({ component: Component, ...rest }) {
  return (
    <Route
      {...rest}
      render={props =>
        RAuth.isAuthenticated ? (
          <Component {...props} />
         
        ) : (<div> <h3> You are not logged in </h3></div>)
      }

    />
  );
}

class  Home extends NavLoggedIn {
  constructor(props) {
    super(props);
    this.state = {
      lastsent : '',
      mqlastsent :'',
      modlastaquired:'',
      statuslog: [],
      tosend : '',

    };
    this.commonChange = this.commonChange.bind(this)
    // this.send2sock = this.send2sock.bind(this)
    // this.send2sockLastsent = this.send2sockLastsent.bind(this)
    // this.send2sockStatuslog= this.send2sockStatuslog.bind(this)
    this.send2sockCommands= this.send2sockCommands.bind(this)
  }
  commonChange(event) {
    // console.log(this.state.tosend)
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  
  // for testing purposes: sending to the echo service which will send it back back
  send2sockCommands (cmd) { 
    // console.log(this.state.tosend)
    WsClient.WsSend(cmd)
  }
  
  componentWillMount(){
    // console.log(WsClient.conn)
    // WsClient.WsConnect()
    WsClient.conn.onmessage = evt => { 
      // console.log("In Home WS Event", evt)
      let rec = JSON.parse(evt.data)
      console.log("In Home WS Event", rec)
      if (rec.hasOwnProperty("mqlastsent")){
        this.setState({
          	mqlastsent : rec.mqlastsent
          })
      } else if (rec.hasOwnProperty("modlastaquired")){
        this.setState({
          modlastaquired : rec.modlastaquired
        })
      }else if  (rec.hasOwnProperty("statuslog")) {
        this.setState({
          statuslog :this.state.statuslog.concat([ rec.statuslog ]).slice(-100)
        })
   
      }
    };
    
  }

  componentWillUnmout() {
    console.log("Unmounting")
    this.connection.close();
  }
  render () { 
  return  <div class="container-fluid ">
          {RAuth.isAuthenticated ? 
            <div> 
                Gateway Process Command
                <input id="send"type="submit" value="Start" onClick = {() => WsClient.WsSendCmd("1") }/>
                <input id="send"type="submit" value="Stop" onClick = {() => WsClient.WsSendCmd("2")}/>
            </div> 
            : null }
          <h3>Last Sent MQTT Message</h3>
          <pre class="box" id="mqtt_msg"> 
          {this.state.mqlastsent}
          </pre>
          <h3>Last Aquired Modbus Message</h3>
          <pre class="box" id="mqtt_msg"> 
          {this.state.modlastaquired}
          </pre>
          <h3>Status </h3>
          
          <ul>{ this.state.statuslog.slice(-100).map( (msg, idx) => <li key={'msg-' + idx }>{ msg}</li> )}</ul>;
        
        </div>
        }
}
            

class  UserManager extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      username :'',
      name : '',
      password : '' ,
      createStatus : ''
    };
    this.commonChange = this.commonChange.bind(this)
    this.handleAddUser = this.handleAddUser.bind(this)
    this.getuserdata = this.getuserdata.bind(this)
    // this.handleRemoveUser= this.handleRemoveUser.bind(this)
  }

  commonChange(event) {
    // console.log(this.state.tosend)
    this.setState({
      [event.target.name]: event.target.value
    });
  }
  getuserdata(){
    PostGet.Get("/user/users", "GET")
           .then(response  => response.json())
           .then(json => {this.setState({ data:json.msg})});
  }
  componentDidMount() {
    this.getuserdata();
  };
  handleAddUser(){
    let bod = { name : this.state.name, username : this.state.username, password : this.state.password}
    PostGet.Post("/user/create", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getuserdata();
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getuserdata();
                          }
                          });
  
  };


  handleRemoveUser = (i) => {
    PostGet.Get("/user/delete/"+i, "DELETE")
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getuserdata(); 
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getuserdata(); 
                          }
                          });
                        
  }

  render() { 
    // const { data } = this.state;
    return (
      <div>
        <label>Name&nbsp;</label><input type="text" name="name"  onChange= { this.commonChange }/> &nbsp;  
        <label>User Name &nbsp;</label><input type="text" name="username" onChange= { this.commonChange }/> &nbsp;
        <label>Password&nbsp;</label><input type="password" name="password"  onChange= { this.commonChange }/> &nbsp;
        <button onClick={this.handleAddUser} className="btn btn-primary float-right" >
            Add User
        </button>
      
      <table class="table" >
        <tr> <th>ID </th> <th>Name</th> <th>User Name</th>  </tr>
        <tbody>{this.state.data.map((item, i) => (
                    
                      <tr key={i}>
                      <td >{item.id}</td><td >{item.name}</td><td >{item.username}</td>
                      <td> 
                      <button className="btn btn-danger btn-sm" onClick={this.handleRemoveUser.bind(this, item.id)} >
                          Remove
                      </button>
                      </td>
                      </tr>
                  
                  )) 
                  } 
        </tbody>
      </table>
      </div>
    );
  }
}

class Interface extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      name :'',
      type : 0,
      ipadd : '',
      port : 0,
      comport :'',
      baudrate : 0,
      databits : 0,
      byteorder :0,
      parity :'',
      stopbits:0,
      timeout:0,
    };
    this.commonChange = this.commonChange.bind(this)
    this.AddModParams = this.AddModParams.bind(this)
    this.getmodparams = this.getmodparams.bind(this)
  }
  commonChange(event) {
    // console.log([event.target.name], event.target.value)
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  getmodparams() {
    PostGet.Get("api/v1/interface/getall", "GET")
            .then(response  => response.json())
            .then(json => {this.setState({ data:json.msg})});
  }

  componentDidMount() {
    this.getmodparams();
  };
  handleRemoveSerial  = (i) => {
    PostGet.Get("api/v1/interface/delete/"+i, "DELETE")
    .then(response  => response.json())
    .then(json => {
                   if (json.msg[0].content === "Done") {
                     // console.log(json.msg)
                     // alert("user")
                     this.getmodparams(); 
                   } else {
                     alert(json.msg[0].content)
                     // console.log(json.msg)
                     this.getmodparams(); 
                   }
                   });
                 
  }

  getType = (i) => {
    switch(i){
      case 1:
        return "Modbus Serial RTU"
      case 3:
        return "Modbus TCP/IP"
      default :
        return "Don't Know"
    }
  }

  AddModParams() {
    var bod = {
               name : this.state.name, type : parseInt(this.state.type), ipadd : this.state.ipadd,
               port : parseInt(this.state.port), comport : this.state.comport, 
               baudrate : parseInt(this.state.baudrate) , databits : parseInt(this.state.databits),  
               parity : this.state.parity, stopbits : parseInt(this.state.stopbits), 
               timeout : parseInt(this.state.timeout) , daqrate : parseInt(this.state.daqrate)
              }
    // console.log(bod)
    PostGet.Post("api/v1/interface/create/0", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                            if (json.msg[0].content === "Done") {
                              // console.log(json.msg)
                              // alert("user")
                              this.getmodparams(); 
                            } else {
                              alert(json.msg[0].content)
                              // console.log(json.msg)
                              this.getmodparams(); 
                            }
            });

  }

  render() { 
    
    return ( <div>
       <div class="row">
       <div class="col-xs-2 panel">
        <label>Name</label> <input class="form-control" type="text" value={this.state.name} name="name" onChange = {this.commonChange }/>
        </div>
        <div class="col-xs-2 panel">
        <label>Type</label> 
          <select  class="form-control" name="type" onChange={this.commonChange}>
            <option ></option>
              <option value ="1">Modbus Serial RTU</option>
              <option value ="2">Modbus TCP/IP</option>
          </select>
        </div>
        <div class="col-xs-2 panel">
        <label>IP Address</label> <input class="form-control" type="text" value={this.state.ipadd} name="ipadd" onChange = {this.commonChange }/>
        </div>
        <div class="col-xs-2 panel">
        <label>Port </label> <input class="form-control" type="text" value={this.state.port} name="port" onChange = {this.commonChange }/>
        </div>
         <div class="col-xs-2 panel">
        <label>Serial Port</label> <input class="form-control" type="text" value={this.state.comport} name="comport" onChange = {this.commonChange }/>
        </div>
        <div class="col-xs-1 panel">
        <label>Baud Rate</label> <input class="form-control" type="number" value={this.state.baudrate} name="baudrate" onChange = {this.commonChange }/>
        </div>
        <div class="col-xs-1 panel">
        <label>Data Bits</label> <input class="form-control" type="number"  value={this.state.databits} name="databits" onChange = {this.commonChange }/> 
        </div>
        
        <div class="col-xs-1 panel">
        <label>Paraity</label> <input class="form-control" type="text"  value={this.state.parity} name="parity" onChange = {this.commonChange }/>
        </div>
        <div class="col-xs-1 panel">
        <label>Stop Bits</label> <input class="form-control"type="number" value={this.state.stopbits} name="stopbits" onChange = {this.commonChange }/> 
        </div>
        <div class="col-xs-1 panel">
        <label>Timeout </label> <input class="form-control" type="number"  value={this.state.timeout} name="timeout" onChange = {this.commonChange }/>
        </div>
        <div class="col-xs-1 panel">
        <label>DAQ Rate </label> <input class="form-control" type="number"  value={this.state.daqrate} name="daqrate" onChange = {this.commonChange }/>
        </div>

        <button onClick={this.AddModParams} className="btn btn-primary pull-bottom" >
          Add Interface
        </button>
        </div>
        <br/>

        <table class="table" >
        <tr> 
          <th>ID </th> <th>Name</th> <th>Type</th> <th>Ip Address</th> 
          <th>Port</th> <th>Serial Port</th> <th>Baud Rate</th> 
          <th>Data Bits</th> <th>Paraity</th> <th>Stop Bits </th> <th>Timeout</th>       
          <th>DAQ Rate</th>     
        </tr>

        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td> <td>{item.name}</td> <td>{this.getType(item.type)}</td>
          <td> {item.ipadd} </td> <td> {item.port} </td>
          <td >{item.comport}</td><td >{item.baudrate}</td>
          <td >{item.databits}</td> <td >{item.parity}</td>
          <td >{item.stopbits}</td> <td >{item.timeout}</td>
          <td >{item.daqrate}</td>
          <td>
          <button className="btn btn-danger btn-sm" onClick={this.handleRemoveSerial.bind(this, item.id)} >
              Remove
          </button>
          </td>
          </tr>        
          )) 
          } 
        </tbody>
      </table>
    </div>);

    }
}

class ModbusRegsiters extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      modparams : {} ,
      name : '',
      functcode : 0,
      register : 0 ,
      qty:0,
      datatype:0,
      byteorder:0,
      postprocess:'',
      tags:'',
      devicemodelsid:0,
    };
    this.commonChange = this.commonChange.bind(this)
    this.handleAddReg = this.handleAddReg.bind(this)
    this.getmodregs = this.getmodregs.bind(this)
   
  
  }

  commonChange(event) {
    // console.log([event.target.name], event.target.value)
    this.setState({
      [event.target.name]: event.target.value
    });
  }
  getmodregs(){
    PostGet.Get("api/v1/modregs/getall", "GET")
           .then(response  => response.json())
           .then(json => {this.setState({ data:json.msg})});
  }
  componentDidMount() {
    this.getmodregs();
  };

  
  
  handleAddReg(){
    let bod = { name : this.state.name ,  functcode : parseInt(this.state.functcode), register : parseInt(this.state.register), qty:parseInt(this.state.qty), datatype:parseInt(this.state.datatype),byteorder:parseInt(this.state.byteorder), postprocess:this.state.postprocess, tags:this.state.tags, devicemodelsid:parseInt(this.state.devicemodelsid),}
    PostGet.Post("api/v1/modregs/create/0", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getmodregs();
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getmodregs();
                          }
                          });
  
  };


  handleRemoveReg = (i) => {
      PostGet.Get("api/v1/modregs/delete/"+i, "DELETE")
      .then(response  => response.json())
      .then(json => {
          if (json.msg[0].content === "Done") {
            // console.log(json.msg)
            // alert("user")
            this.getmodregs(); 
          } else {
            alert(json.msg[0].content)
            // console.log(json.msg)
            this.getmodregs(); 
          }
          });
}
  getDataTypeName = (i) => {
    switch(i){
      case 1:
        return "Uint8"
      case 2:
        return "Uint8 array"
      case 3:
        return "Int8"
      case 4:
        return "Int8 Array"
        case 5:
        return "Uint16"
      case 6:
        return "Uint16 Array"
      case 7:
        return "Int16"
      case 8:
        return "Int16 Array"
        case 9:
        return "Uint32"
      case 10:
        return "Uint32 Array"
      case 11:
        return "Int32"
      case 12:
        return "Int32 Array"
        case 13:
        return "Uint64"
      case 14:
        return "Uint64 Array"
      case 15:
        return "Int64"
      case 16:
        return "Int64 Array"
        case 17:
        return "Float32"
      case 18:
        return "Float32 Array"
      case 19:
        return "Float64"
      case 20:
        return "Float64 Array"
       
      default:
        return "I dont know"
    }
  }
  getFunctionCode = (i) => {
    switch(i){
      case 1:
        return "Read Coil FC-1" 
      case 2:
        return "Read Discrete Input FC-2"
      case 3:
        return "Read Holding Registers FC-3"
      case 4:
        return "Read Input Registers FC-4"
      case 5:
        return "Write Single Coil FC-5" 
      case 6:
        return "Write Single Holding Register FC-6"
      case 15:
        return "Write Multiple Coils FC-15"
      case 16:
        return "Write Multiple Holding Registers FC-16"
      default:
        return "I dont know"
  }
  }
  
  getByteOrderText = (i) => {
    switch(i){
      case 1:
        return "Big Endian (MSB->ABCD<-LSB)"
      case 2:
        return "Little Endian (MSB->DCBA<-LSB)"
      case 3:
        return "Mid-Big Endian (MSB->BADC<-LSB)"
      case 4:
        return "Mid-Little Endian (MSB->CDAB<-LSB)"
      default:
        return "I dont know"
  }
  }
  render() { 
    return (
      <div>
        <div class="col-xs-2 panel">
        <label>Name &nbsp;</label><input class="form-control" type="text" name="name"  onChange= { this.commonChange }/> &nbsp;
        </div>

        

        <div class="col-xs-2 panel">
        <label>Function Code &nbsp;</label>
        <select  class="form-control" name="functcode" onChange={this.commonChange}>
          <option ></option>
            <option value ="1">Read Coil FC-1</option>
            <option value ="2">Read Discrete Input FC-2</option>
            <option value ="3">Read Holding Registers FC-3</option>
            <option value ="4">Read Input Registers FC-4</option>
          

        </select>
        &nbsp;
        </div>

        <div class="col-xs-1 panel">
        <label>Register&nbsp;</label><input class="form-control" type="number" name="register"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-1 panel">
        <label>qty&nbsp;</label><input class="form-control" type="number" name="qty"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-2 panel">
        <label>DataType&nbsp;</label>
        <select  class="form-control" name="datatype" onChange={this.commonChange}>
          <option ></option>
            <option value ="1">Uint8</option>
            <option value ="2">Uint8 Array</option>
            <option value ="3">Int8</option>
            <option value ="4">Int8 Array</option>
            <option value ="5">Uint16</option>
            <option value ="6">Uint16 Array</option>
            <option value ="7">Int16</option>
            <option value ="8">Int16 Array</option>
            <option value ="9">Uint32</option>
            <option value ="10">Uint32 Array</option>
            <option value ="11">Int32</option>
            <option value ="12">Int32 Array</option>
            <option value ="13">Uint64</option>
            <option value ="14">Uint64 Array</option>
            <option value ="15">Int64</option>
            <option value ="16">Int64 Array</option>
            <option value ="17">Float32</option>
            <option value ="18">Float32 Array</option>
            <option value ="19">Float64</option>
            <option value ="20">Float64 Array</option>            
        </select>
        </div>

        &nbsp;

        <div class="col-xs-2 panel">
        <label>Byte Order(Endian)&nbsp;</label>
        <select  class="form-control" name="byteorder" onChange={this.commonChange}>
          <option ></option>
            <option value ="1">Big Endian {"(MSB->ABCD<-LSB)"}</option>
            <option value ="2">Little Endian {"(MSB->DCBA<-LSB)"}</option>
            <option value ="3">Mid-Big Endian {"(MSB->BADC<-LSB)"}</option>
            <option value ="4">Mid-Little Endian {"(MSB->CDAB<-LSB)"}</option>
        </select>
        </div>

        <div class="col-xs-2 panel">
        <label>PostProcess&nbsp;</label><input class="form-control" type="text" name="postprocess"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-2 panel">
        <label>Tags&nbsp;</label><input class="form-control" type="text" name="tags"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-2 panel">
        <label>Device Model ID &nbsp;</label><input class="form-control" type="number" name="devicemodelsid"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-1 panel">
          <button onClick={this.handleAddReg} className="btn btn-primary float-right" >
              Add Register
          </button>
        </div>

        


      <table class="table" >
        <tr> 
          <th>ID </th> <th>Name</th>  <th>Function Code</th>
          <th>Register</th> <th>Qty</th> <th>DataType</th> <th>Byte Order (Endian)</th>
          <th>PostProcess</th><th>Tags</th> <th>Model Id</th>
          <th> Remove </th>
        </tr>
        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td> <td >{item.name}</td>
          <td >{this.getFunctionCode(item.functcode)}</td>
          <td >{item.register}</td> <td >{item.qty}</td>
          <td >{this.getDataTypeName(item.datatype)}</td> 
          <td >{this.getByteOrderText(item.byteorder)}</td> 
          <td >{item.postprocess}</td>
          <td >{item.tags.split(",").map( (tag,i) => ( <tr key={i}> <td> {tag} </td></tr> ) )}</td>
          <td> {item.devicemodelsid}</td>
          <td>
          <button className="btn btn-danger btn-sm" onClick={this.handleRemoveReg.bind(this, item.id)} >
              Remove
          </button>
          </td>
          </tr>        
          )) 
          } 
        </tbody>
      </table>
      </div>
    );
  }
}


class Mqtt extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      mqparams : {},
      ip : '',
      port : 0 ,
      username : '',
      password : '',

      
    };
    this.commonChange = this.commonChange.bind(this)
    this.getmqparams = this.getmqparams.bind(this)
  
  }

  commonChange(event) {
    // console.log( [event.target.name], event.target.value)  
    this.setState({
      [event.target.name]: event.target.value
    });
  }
  
  getmqparams(){ 
    PostGet.Get("api/v1/mqtt/getall", "GET")
           .then(response  => response.json())
           .then(json => {
                          this.setState({ mqparams:json.msg[0], ip:json.msg[0].ip, port:json.msg[0].port, username:json.msg[0].username, password:json.msg[0].password })
                         }
                );
  }

  UpdateMqParams = () => {
    var bod = { ip : this.state.ip, port : parseInt(this.state.port) , username : this.state.username, password : this.state.password}
    // console.log(bod)
    PostGet.Post("api/v1/mqtt/create/1", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                            if (json.msg[0].content === "Done") {
                              // console.log(json.msg)
                              // alert("user")
                              this.getmqparams(); 
                            } else {
                              alert(json.msg[0].content)
                              // console.log(json.msg)
                              this.getmqparams(); 
                            }
            });
  }



  componentDidMount() {
    this.getmqparams();
  }


  render() { 
    return (
      <div>
        <div class ="row"> 
        <label>&nbsp;&nbsp;&nbsp; Ip Or Hostname</label> 
        <input type="text" name='ip' value={this.state.ip} onChange={this.commonChange} /> &nbsp;
        <label>Port</label> <input   type="text" name="port" value = {this.state.port} onChange= { this.commonChange }/> &nbsp;  
        <label>User Name </label> <input  id="mqtt_username" type="text" name="username" value = {this.state.username} onChange= { this.commonChange } readOnly = {this.state.readonly}  onClick= { () => { this.setState({readonly : false})} }/> &nbsp;  
        <label>Password</label> <input id="mqtt_password"  type="password" name="password" value = {this.state.password} onChange= { this.commonChange }  readOnly = {this.state.readonly} onClick = { () => { this.setState({readonly : false})} }  /> &nbsp;  
        <button onClick={this.UpdateMqParams} className="btn btn-primary float-right" >
            Update
        </button>
        </div>
      </div>
    );
  }
}


class DeviceModels extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      make : '',
      model:'',
    };
    this.commonChange = this.commonChange.bind(this)
    this.getdevicemodels = this.getdevicemodels.bind(this)
  }

  commonChange(event) {
    // console.log([event.target.name], event.target.value)
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  getdevicemodels(){
    PostGet.Get("api/v1/devicemodels/getall", "GET")
           .then(response  => response.json())
           .then(json => { this.setState({ data:json.msg})});
  }

  componentDidMount() {
    this.getdevicemodels();
  };

  
  
  handleAddDevModel = (i) => {
    let bod = { make : this.state.make , model : this.state.model,}
    PostGet.Post("api/v1/devicemodels/create/0", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getdevicemodels();
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getdevicemodels();
                          }
                          });
  
  }


  handleRemoveDevModel = (i) => {
      PostGet.Get("api/v1/devicemodels/delete/"+i, "DELETE")
      .then(response  => response.json())
      .then(json => {
          if (json.msg[0].content === "Done") {
            // console.log(json.msg)
            // alert("user")
            this.getdevicemodels(); 
          } else {
            alert(json.msg[0].content)
            // console.log(json.msg)
            this.getdevicemodels(); 
          }
          });
}
  getDataTypeName = (i) => {
    switch(i){
      case 1:
        return "Uint8"
      case 2:
        return "Uint8 array"
      case 3:
        return "Int8"
      case 4:
        return "Int8 Array"
        case 5:
        return "Uint16"
      case 6:
        return "Uint16 Array"
      case 7:
        return "Int16"
      case 8:
        return "Int16 Array"
        case 9:
        return "Uint32"
      case 10:
        return "Uint32 Array"
      case 11:
        return "Int32"
      case 12:
        return "Int32 Array"
        case 13:
        return "Uint64"
      case 14:
        return "Uint64 Array"
      case 15:
        return "Int64"
      case 16:
        return "Int64 Array"
        case 17:
        return "Float32"
      case 18:
        return "Float32 Array"
      case 19:
        return "Float64"
      case 20:
        return "Float64 Array"
       
      default:
        return "I dont know"
    }
  }
  getFunctionCode = (i) => {
    switch(i){
      case 1:
        return "Read Coil FC-1" 
      case 2:
        return "Read Discrete Input FC-2"
      case 3:
        return "Read Holding Registers FC-3"
      case 4:
        return "Read Input Registers FC-4"
      case 5:
        return "Write Single Coil FC-5" 
      case 6:
        return "Write Single Holding Register FC-6"
      case 15:
        return "Write Multiple Coils FC-15"
      case 16:
        return "Write Multiple Holding Registers FC-16"
      default:
        return "I dont know"
  }
  }
  
  getByteOrderText = (i) => {
    switch(i){
      case 1:
        return "Big Endian (MSB->ABCD<-LSB)"
      case 2:
        return "Little Endian (MSB->DCBA<-LSB)"
      case 3:
        return "Mid-Big Endian (MSB->BADC<-LSB)"
      case 4:
        return "Mid-Little Endian (MSB->CDAB<-LSB)"
      default:
        return "I dont know"
  }
  }
  render() { 
    return (
      <div>
        <div class="col-xs-2 panel">
        <label>Make &nbsp;</label><input class="form-control" type="text" name="make"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-1 panel">
        <label>Model &nbsp;</label><input class="form-control" type="text" name="model"  onChange= { this.commonChange }/> &nbsp;  
        </div>


        <div class="col-xs-1 panel">
          <button onClick={this.handleAddDevModel} className="btn btn-primary float-right" >
              Add Register
          </button>
        </div>

        


      <table class="table" >
        <tr> 
          <th>ID </th> <th>Make</th>  <th>Model</th>
          <th> Remove </th>
        </tr>
        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td> <td >{item.make}</td>
          <td >{item.model}</td>
          <td>
          <button className="btn btn-danger btn-sm" onClick={this.handleRemoveDevModel.bind(this, item.id)} >
              Remove
          </button>
          </td>
          </tr>        
          )) 
          } 
        </tbody>
      </table>
      </div>
    );
  }
}

class Devices extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      name : '',
      deviceid:'',
      mbid:0,
      devicemodelsid : 0,
      intefacedetailsid: 0,
    };
    this.commonChange = this.commonChange.bind(this)
    this.getdevices = this.getdevices.bind(this)
  }

  commonChange(event) {
    // console.log([event.target.name], event.target.value)
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  getType = (i) => {
    switch(i){
      case 1:
        return "Modbus RTU"
      case 2:
        return "ModbusASCII"
      case 3:
        return "Modbus TCP"
      default:
        return "I dont know"
  }
}

  getdevices(){
    PostGet.Get("api/v1/devices/getall", "GET")
           .then(response  => response.json())
           .then(json => {this.setState({ data:json.msg})});
  }

  componentDidMount() {
    this.getdevices();
  };

  
  
  handleAddDevices = (i) => {
    let bod = { name : this.state.name , deviceid : this.state.deviceid, mbid : parseInt(this.state.mbid), devicemodelsid : parseInt(this.state.devicemodelsid), intefacedetailsid:parseInt(this.state.intefacedetailsid)}
    // console.log(bod)
    PostGet.Post("api/v1/devices/create/0", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getdevices();
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getdevices();
                          }
                          });
  
  }


  handleRemoveDevices = (i) => {
      PostGet.Get("api/v1/devices/delete/"+i, "DELETE")
      .then(response  => response.json())
      .then(json => {
          if (json.msg[0].content === "Done") {
            // console.log(json.msg)
            // alert("user")
            this.getdevices(); 
          } else {
            alert(json.msg[0].content)
            // console.log(json.msg)
            this.getdevices(); 
          }
          });
}

 
  render() { 
    return (
      <div>
        <div class="col-xs-2 panel">
        <label>Name &nbsp;</label><input class="form-control" type="text" name="name"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-2 panel">
        <label>Device ID &nbsp;</label><input class="form-control" type="text" name="deviceid"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-2 panel">
        <label>Modbus ID &nbsp;</label><input class="form-control" type="number" name="mbid"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-1 panel">
        <label>Device Model ID &nbsp;</label><input class="form-control" type="number" name="devicemodelsid"  onChange= { this.commonChange }/> &nbsp;  
        </div>

        <div class="col-xs-1 panel">
        <label>Interface ID &nbsp;</label><input class="form-control" type="number" name="intefacedetailsid"  onChange= { this.commonChange }/> &nbsp;  
        </div>

        <div class="col-xs-1 panel">
          <button onClick={this.handleAddDevices} className="btn btn-primary float-right" >
              Add Register
          </button>
        </div>
    
        <table class="table" >
          <tr> 
            <th>ID </th> <th>Name</th>  <th>Device ID</th>
            <th>Modbus ID</th>  <th>Make</th><th>Model</th>
            <th>Interface </th><th>Type</th>
            <th> Remove </th>
          </tr>
          <tbody>{this.state.data.map((item, i) => (        
            <tr key={i}> 
            <td >{item.Device.id}</td> <td >{item.Device.name}</td>
            <td >{item.Device.deviceid}</td><td >{item.Device.mbid}</td>
            <td >{item.Model.make}</td> <td >{item.Model.model}</td> 
            <td >{item.Interface.name}</td>
            <td >{this.getType(item.Interface.type)}</td>
            <td>
            <button className="btn btn-danger btn-sm" onClick={this.handleRemoveDevices.bind(this, item.Device.id)} >
                Remove
            </button>
            </td>
            </tr>        
            )) 
            } 
          </tbody>
        </table>
      </div>
    );
  }
}

class AdditionalFeatures extends React.Component{
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      param : '',
      value:'',
      paramtype:0,
    };
    this.commonChangeArray = this.commonChangeArray.bind(this)
    this.getaddfeatures = this.getaddfeatures.bind(this)
  }

  commonChangeArray(event) {
    // console.log( [event.target.id],[event.target.name], event.target.value  ) 
    let newArray = [...this.state.data];
    newArray[event.target.id][event.target.name] = event.target.value
    this.setState({ data : newArray});
  }

  componentDidMount() {
    this.getaddfeatures();
  };
  getaddfeatures(){
    PostGet.Get("api/v1/addfeatures/getall", "GET")
           .then(response  => response.json())
           .then(json => {console.log(json.msg); this.setState({ data:json.msg})});
  }

  handleUpdateAddFeatures = (i) => {
    let bod = this.state.data[i-1] 
    // console.log(bod)
    PostGet.Post("api/v1/addfeatures/create/"+i, "POST",bod)
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getaddfeatures();
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getaddfeatures();
                          }
                          });
    
  }
  handleAddAddFeatures = (i) => {
    let bod = { param : this.state.param , value : this.state.value, paramtype : parseInt(this.state.paramtype), }
    // console.log(bod)
    PostGet.Post("api/v1/addfeatures/create/0", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getaddfeatures();
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getaddfeatures();
                          }
                          });
  }

  handleDelAddFeatures = (i) => {
    PostGet.Get("api/v1/addfeatures/delete/"+i, "DELETE")
      .then(response  => response.json())
      .then(json => {
          if (json.msg[0].content === "Done") {
            // console.log(json.msg)
            // alert("user")
            this.getaddfeatures(); 
          } else {
            alert(json.msg[0].content)
            // console.log(json.msg)
            this.getaddfeatures(); 
          }
          });
  }

  render() { 
    return (
      <div>
        {/* <div class="col-xs-2 panel">
        <label>Parameter Name &nbsp;</label><input class="form-control" type="text" name="param"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-2 panel">
        <label>Value  &nbsp;</label><input class="form-control" type="text" name="value"  onChange= { this.commonChange }/> &nbsp;
        </div>

        <div class="col-xs-2 panel">
        <label>Parameter type :  1 - Bool ,  2 - UInt64, 3 - Int64, 4 - Float64 &nbsp;</label><input class="form-control" type="number" name="paramtype"  onChange= { this.commonChange }/> &nbsp;
        </div>


       
        <div class="col-xs-1 panel">
          <button onClick={this.handleAddDevices} className="btn btn-primary float-right" >
              Add Features
          </button>
        </div> */}
        <label>Parameter type : 0 - string  1 - Bool ,  2 - UInt64, 3 - Int64, 4 - Float64 &nbsp;</label>
        <table class="table" >
          <tr> 
            <th>ID </th> <th>Parameter name</th>  <th>Value</th>
            <th>Parameter Type </th>  
            <th></th><th></th><th></th>
          </tr>
          <tbody>{this.state.data.map((item, i) => (        
            <tr key={i}> 
            <td >{item.id}</td> <td >{item.param}</td>
            <td>
              <input  id={i} type="text" name="value" value = {item.value} onChange= { this.commonChangeArray } readOnly = {this.state.readonly}  onClick= { () => { this.setState({readonly : false})} }/>
            </td>
            <td>
              <input  id={i} type="text" name="paramtype" value = {item.paramtype} onChange= { this.commonChangeArray } readOnly = {this.state.readonly}  onClick= { () => { this.setState({readonly : false})} }/>
            </td>
            <td >{item.value}</td><td >{item.paramtype}</td>
            
            <td>
            <button className="btn btn-success btn-sm" onClick={this.handleUpdateAddFeatures.bind(this, item.id)} >
                Update
            </button>
            </td>
            </tr>        
            )) 
            } 
          </tbody>
        </table>
      </div>
    );
  }
}

          
ReactDOM.render(<AuthExample />, document.getElementById('app'));            