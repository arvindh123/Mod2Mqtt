var { BrowserRouter, Link, withRouter, Route, Redirect } = ReactRouterDOM;
var Router = BrowserRouter;
// var {createReactClass} = ReactRouterDOM;


class AuthExample extends React.Component {
  componentWillMount() {
    WsClient.WsConnect();
  }
  render() {

    return (
      <Router basename="/#">
        <div>
          <AuthButton />
          <Redirect to="/home" />
          <Route path="/home" component={Home} />
          <Route path="/login" component={NavLogIn} />
          <PrivateRoute path="/usermanager" component={UserManager} />
          <PrivateRoute path="/serial" component={SerialPorts} />
          <PrivateRoute path="/modbus" component={Modbus} />
          <PrivateRoute path="/smr" component={SMR} />
          <PrivateRoute path="/mqtt" component={Mqtt} />
          <PrivateRoute path="/trr" component={TRR} />

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
      // console.log(response);
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
              <Link to="/serial">Interface</Link>
            </li>
            <li>
              <Link to="/modbus">Modbus</Link>
            </li>
            <li>
              <Link to="/smr">Serail ModRegs Relation</Link>
            </li>
            <li>
              <Link to="/mqtt">Mqtt</Link>
            </li>
            <li>
              <Link to="/trr">Topics and Registers Relation</Link>
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
              <Link to="/serial">Interface</Link>
            </li>
            <li>
              <Link to="/modbus">Modbus</Link>
            </li>
          
            <li>
              <Link to="/mqtt">Mqtt</Link>
            </li>
            <li>
              <Link to="/trr">Topics and Registers Relation</Link>
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
class SerialPorts extends React.Component {
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
    PostGet.Get("api/v1/modbus/params", "GET")
            .then(response  => response.json())
            .then(json => {this.setState({ data:json.msg})});
  }

  componentDidMount() {
    this.getmodparams();
  };
  handleRemoveSerial  = (i) => {
    PostGet.Get("api/v1/modbus/params/"+i, "DELETE")
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
      case 2:
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
               timeout : parseInt(this.state.timeout) 
              }
    // console.log(bod)
    PostGet.Post("api/v1/modbus/params/0", "POST",bod)
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

        <button onClick={this.AddModParams} className="btn btn-primary pull-bottom" >
          Add Serial
        </button>
        </div>
        <br/>

        <table class="table" >
        <tr> 
          <th>ID </th> <th>Name</th> <th>Type</th> <th>Ip Address</th> 
          <th>Port</th> <th>Serial Port</th> <th>Baud Rate</th> 
          <th>Data Bits</th> <th>Paraity</th> <th>Stop Bits </th> <th>Timeout</th>         
        </tr>

        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td> <td>{item.name}</td> <td>{this.getType(item.type)}</td>
          <td> {item.ipadd} </td> <td> {item.port} </td>
          <td >{item.comport}</td><td >{item.baudrate}</td>
          <td >{item.databits}</td> <td >{item.parity}</td>
          <td >{item.stopbits}</td> <td >{item.timeout}</td>
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

class Modbus extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      modparams : {} ,
      name : '',
      unit :0,
      functcode : 0,
      register : 0 ,
      qty:0,
      datatype:0,
      postprocess:'',
      tags:'',
      
      
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
    PostGet.Get("api/v1/modbus/regs", "GET")
           .then(response  => response.json())
           .then(json => {this.setState({ data:json.msg})});
  }
  componentDidMount() {
    this.getmodregs();
  };

  
  
  handleAddReg(){
    let bod = { name : this.state.name , unit : parseInt(this.state.unit), functcode : parseInt(this.state.functcode), register : parseInt(this.state.register), qty:parseInt(this.state.qty), datatype:parseInt(this.state.datatype),byteorder:parseInt(this.state.byteorder), postprocess:this.state.postprocess, tags:this.state.tags,}
    PostGet.Post("api/v1/modbus/regs/0", "POST",bod)
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
    PostGet.Get("api/v1/modbus/regs/"+i, "DELETE")
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
      case 90:
        return "Suid"
      case 99:
        return "Unix Ts"
      case 98:
        return "Ts with UTC"
      case 97:
        return "Ts with Format in Post Process"
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

        <div class="col-xs-1 panel">
        <label>Unit &nbsp;</label><input class="form-control" type="number" name="unit"  onChange= { this.commonChange }/> &nbsp;  
        </div>

        <div class="col-xs-2 panel">
        <label>Function Code &nbsp;</label>
        <select  class="form-control" name="functcode" onChange={this.commonChange}>
          <option ></option>
            <option value ="1">Read Coil FC-1</option>
            <option value ="2">Read Discrete Input FC-2</option>
            <option value ="3">Read Holding Registers FC-3</option>
            <option value ="4">Read Input Registers FC-4</option>
            <option value ="90">Suid</option>
            <option value ="97">Ts with Format in Post Process</option>
            <option value ="98">Ts with UTC</option>
            <option value ="99">Unix Ts</option>

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

        <div class="col-xs-1 panel">
          <button onClick={this.handleAddReg} className="btn btn-primary float-right" >
              Add Register
          </button>
        </div>


      <table class="table" >
        <tr> 
          <th>ID </th> <th>Name</th> <th>Unit</th> <th>Function Code</th>
          <th>Register</th> <th>Qty</th> <th>DataType</th> <th>Byte Order (Endian)</th>
          <th>PostProcess</th><th>Tags</th> <th>Interface Name</th>
          {/* <th>Topic ID </th>  */}
          <th>Topic</th> 
          {/* <th>Qos</th> <th>Retain</th> */}
          <th> Remove </th>
          {/* <th> </th><th> </th><th> </th> */}
        </tr>
        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td> <td >{item.name}</td>
          <td >{item.unit}</td> <td >{this.getFunctionCode(item.functcode)}</td>
          <td >{item.register}</td> <td >{item.qty}</td>
          <td >{this.getDataTypeName(item.datatype)}</td> 
          <td >{this.getByteOrderText(item.byteorder)}</td> 
          <td >{item.postprocess}</td>
          <td >{item.tags.split(",").map( (tag,i) => ( <tr key={i}> <td> {tag} </td></tr> ) )}</td>
          <td>{item.serialport.map((item, i) => (
             <tr key={i}> <td >{item.name}</td> </tr>)) }
      
            </td>
          <td> {item.mqtopic.map((item, i) => (
             <tr key={i}>
            {/* <td >{item.id}</td><td >{item.topic}</td><td >{item.qos}</td><td >{String(item.retain)}</td> */}
            <td >{item.topic}</td>
            </tr>
          )) }
          </td>
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

class SMR extends React.Component {
  // SMR - Serial Port and  Modbus Registers  Relationship mapping
  constructor(props) {
    super(props);
    this.state = {
      options : [ ],
      selectedRegs : [],
      ports : [],
      selectedPort : 0
    };
    this.getmodregs = this.getmodregs.bind(this)
    this.getports = this.getports.bind(this)
    this.handlePortChange = this.handlePortChange.bind(this)
  }

  handleDeselect = (deselectedOptions) => {
    var selectedRegs = this.state.selectedRegs.slice()
    deselectedOptions.forEach(option => {
      selectedRegs.splice(selectedRegs.indexOf(option), 1)
    })
    this.setState({selectedRegs})
  }

  handleSelectionChange = (selectedRegs) => {
    this.setState({selectedRegs})
  }
  
  getports(){
    PostGet.Get("api/v1/modbus/params", "GET")
           .then(response  => response.json())
           .then(json => {  //console.log(json.msg);
                            this.setState({ports : json.msg }) ;
                            if(json.msg[0].hasOwnProperty("modregs")){
                            this.setState ({ selectedRegs  :   json.msg[0].modregs})
                            }
            });
  }

  getmodregs(){
    PostGet.Get("api/v1/modbus/regs", "GET")
           .then(response  => response.json())
           .then(json => { this.setState({options : json.msg })});
  }

  componentDidMount(){
    this.getports()
    this.getmodregs()
    
  }

  handlePortChange(event) {
    this.setState({selectedPort:event.target.value})
    this.setState ( { selectedRegs  : this.state.ports[event.target.value].modregs})
    // console.log(event.target.value)
    
  }

  SaveTheRelation = () =>  {

    var modregids = [] 
    var selectedRegs = this.state.selectedRegs
    selectedRegs.forEach(reg => 
      modregids.push(reg.id)
      ) 
      var bod = { modregids :modregids }
      var portid = this.state.ports[this.state.selectedPort].id

    PostGet.Get("api/v1/serial/modregs/" + portid + "/all","DELETE")
           .then(response  => response.json())
           .then(json => { 
                if (json.msg[0].content === "Done") {
                  if (modregids.length > 0) { 
                    var bod = { modregids :modregids }
                    PostGet.Post("api/v1/serial/modregs/" + portid, "POST",bod)
                          .then(response => response.json())
                          .then(json => {
                            if (json.msg[0].content === "Done" ) {
                              this.getports()
                            }else{
                              alert(json.msg[0].content)
                            }
                          })
                  }
                }else {
                  alert(json.msg[0].content)
                }
           });
  }

  render() { 
    return (
        <div >
          <div class="row">
          <div className="col-md-6">
          <label>Select the Interface </label> 
          <select class="form-control" value={this.state.selectedPort} onChange={this.handlePortChange}>
          {this.state.ports.map((item,idx) => (
            <option key= {idx} value={idx}>{ item.name}</option>
          ))}
          </select> 
          </div>
          <div className="col-md-6">  
          <button class="btn btn-success pull-center" type="button" onClick={() => this.SaveTheRelation()}> Save </button>
          </div>
          </div>
          <div className="col-md-6">
            <h4>Add Modbus Register to Interface </h4>
            <FilteredMultiSelect buttonText= "Add to Interface" className='row' classNames ={{ button:"btn btn-primary", buttonActive:"btn btn-success", filter:"form-control", select:"form-control"}} onChange={this.handleSelectionChange} textProp="name"  valueProp="id" options={this.state.options} selectedOptions={this.state.selectedRegs} />
            </div>
            <div className="col-md-6">
            <h4>Remove Modbus Register from Interface </h4>
            <FilteredMultiSelect buttonText= "Remove from Interface" className='row' classNames ={{ button:"btn btn-primary", buttonActive:"btn btn-danger", filter:"form-control", select:"form-control"}} onChange={this.handleDeselect} textProp="name"  valueProp="id" options={this.state.selectedRegs}  />
            </div>       
        </div>
    )
  }
}

class Mqtt extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      mqparams : {},
      topic :0,
      qos : 0,
      retain : "false" ,
      delay : 0,
      ip : '',
      port : 0 ,
      username : '',
      password : '',
      readonly : true,
      
    };
    this.commonChange = this.commonChange.bind(this)
    this.handleAddMqTopic = this.handleAddMqTopic.bind(this)
    this.getmqtopics = this.getmqtopics.bind(this)
    this.getmqparams = this.getmqparams.bind(this)
  
  }

  commonChange(event) {
    // console.log( [event.target.name], event.target.value)  
    this.setState({
      [event.target.name]: event.target.value
    });
  }
  
  getmqparams(){ 
    PostGet.Get("api/v1/mqtt/params", "GET")
           .then(response  => response.json())
           .then(json => {
                          this.setState({ mqparams:json.msg[0], ip:json.msg[0].ip, port:json.msg[0].port, username:json.msg[0].username, password:json.msg[0].password })
                         }
                );
  }

  UpdateMqParams = () => {
    var bod = { ip : this.state.ip, port : parseInt(this.state.port) , username : this.state.username, password : this.state.password}
    // console.log(bod)
    PostGet.Post("api/v1/mqtt/params", "POST",bod)
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

  getmqtopics(){
    PostGet.Get("api/v1/mqtt/topics", "GET")
           .then(response  => response.json())
           .then(json => {this.setState({ data:json.msg})});
  }

  componentDidMount() {
    this.getmqtopics();
    this.getmqparams();
  }

  handleAddMqTopic(){
    let bod = { topic : this.state.topic, qos : parseInt(this.state.qos), retain : this.state.retain.toLowerCase() == 'true' ? true : false, delay :parseInt(this.state.delay) }
    // console.log(bod)
    PostGet.Post("api/v1/mqtt/topic/0", "POST",bod)
           .then(response  => response.json())
           .then(json => {
                          if (json.msg[0].content === "Done") {
                            // console.log(json.msg)
                            // alert("user")
                            this.getmqtopics();
                          } else {
                            alert(json.msg[0].content)
                            // console.log(json.msg)
                            this.getmqtopics();
                          }
                          });
  
  };


  handleRemovemqtopics = (i) => {
    PostGet.Get("api/v1/mqtt/topics/"+i, "DELETE")
           .then(response  => response.json())
           .then(json => {
                            if (json.msg[0].content === "Done") {
                              // console.log(json.msg)
                              // alert("user")
                              // setTimeout(function(){}, 1000);
                              this.getmqtopics(); 
                            } else {
                              alert(json.msg[0].content)
                              // console.log(json.msg)
                              this.getmqtopics(); 
                            }
                          });
                        
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
        <label>Topic &nbsp;</label><input type="text" name="topic"  onChange= { this.commonChange }/> &nbsp;  
        <label>QoS &nbsp;</label><input type="number" name="qos" onChange= { this.commonChange }/> &nbsp;
        <label>Retain&nbsp;</label><input type="text" name="retain"  onChange= { this.commonChange }/> &nbsp;
        <label>Delay For Visulaization &nbsp;</label><input type="text" name="delay"  onChange= { this.commonChange }/> &nbsp;
        <button onClick={this.handleAddMqTopic} className="btn btn-primary float-right" >
            Add Topic
        </button>
      <table class="table" >
        <tr> <th>ID </th> <th>Topic</th> <th>QoS</th> <th>Retain</th> <th>Delay</th> </tr>
        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td><td >{item.topic}</td><td >{item.qos}</td><td >{String(item.retain)}</td><td >{item.delay}</td>
          <td> 
          <table class="table">
          <tr> <th>ID </th> <th>Name</th> <th>Unit</th> <th>Function Code</th> <th>Register</th> <th>Qty</th> </tr>
          <tbody>{item.modregs.map((item, i) => (
             <tr key={i}>
            <td >{item.id}</td><td >{item.name}</td><td >{item.unit}</td><td >{item.functcode}</td><td >{item.register}</td><td >{item.qty}</td>
            </tr>
          )) }
          </tbody>
          </table>
          </td>
          <td>
          <button className="btn btn-danger btn-sm" onClick={this.handleRemovemqtopics.bind(this, item.id)} >
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

class TRR extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      options : [ ],
      selectedRegs : [],
      topics : [],
      selectedTopic : 0
    };
    this.getmodregs = this.getmodregs.bind(this)
    this.getmqtopics = this.getmqtopics.bind(this)
    this.handleTopicChange = this.handleTopicChange.bind(this)
  }

  handleDeselect = (deselectedOptions) => {
    var selectedRegs = this.state.selectedRegs.slice()
    deselectedOptions.forEach(option => {
      selectedRegs.splice(selectedRegs.indexOf(option), 1)
    })
    this.setState({selectedRegs})
  }

  handleSelectionChange = (selectedRegs) => {
    this.setState({selectedRegs})
  }
  
  getmqtopics(){
    PostGet.Get("api/v1/mqtt/topics", "GET")
           .then(response  => response.json())
           .then(json => {  //console.log(json.msg);
                            this.setState({topics : json.msg }) ;
                            if(json.msg[0].hasOwnProperty("modregs")){
                            this.setState ({ selectedRegs  :   json.msg[0].modregs})
                            }
          
            });
  }

  getmodregs(){
    PostGet.Get("api/v1/modbus/regs", "GET")
           .then(response  => response.json())
           .then(json => { this.setState({options : json.msg })});
  }

  componentDidMount(){
    this.getmqtopics()
    this.getmodregs()
    
  }

  handleTopicChange(event) {
    this.setState({selectedTopic:event.target.value})
    this.setState ( { selectedRegs  : this.state.topics[event.target.value].modregs})
    // console.log(event.target.value)
    
  }

  SaveTheRelation = () =>  {
    // console.log("selected regs .." , this.state.selectedRegs)
    // console.log("selected topic .." ,this.state.selectedTopic)
    var modregids = [] 
    var selectedRegs = this.state.selectedRegs
    selectedRegs.forEach(reg => 
      modregids.push(reg.id)
      )
      // console.log(modregids) 
      var bod = { modregids :modregids }
      var topicid = this.state.topics[this.state.selectedTopic].id
      // console.log(bod)

    PostGet.Get("api/v1/topics/modregs/" + topicid + "/all","DELETE")
           .then(response  => response.json())
           .then(json => { 
                if (json.msg[0].content === "Done") {
                  if (modregids.length > 0) { 
                    var bod = { modregids :modregids }
                    PostGet.Post("api/v1/topics/modregs/" + topicid, "POST",bod)
                          .then(response => response.json())
                          .then(json => {
                            if (json.msg[0].content === "Done" ) {
                              this.getmqtopics()
                            }else{
                              alert(json.msg[0].content)
                            }
                          })
                  }
                }else {
                  alert(json.msg[0].content)
                }
           });
  }

  render() { 
    return (
        <div >
          <div class="row">
          <div className="col-md-6">
          <label>Select the Topic </label> 
          <select class="form-control" value={this.state.selectedTopic} onChange={this.handleTopicChange}>
          {this.state.topics.map((item,idx) => (
            <option key= {idx} value={idx}>{ item.topic}</option>
          ))}
          </select> 
          </div>
          <div className="col-md-6">  
          <button class="btn btn-success pull-center" type="button" onClick={() => this.SaveTheRelation()}> Save </button>
          </div>
          </div>
          <div className="col-md-6">
            <h4>Add Modbus Register to Topics </h4>
            <FilteredMultiSelect buttonText= "Add to Topic" className='row' classNames ={{ button:"btn btn-primary", buttonActive:"btn btn-success", filter:"form-control", select:"form-control"}} onChange={this.handleSelectionChange} textProp="name"  valueProp="id" options={this.state.options} selectedOptions={this.state.selectedRegs} />
            </div>
            <div className="col-md-6">
            <h4>Remove Modbus Register from Topics </h4>
            <FilteredMultiSelect buttonText= "Remove from Topic" className='row' classNames ={{ button:"btn btn-primary", buttonActive:"btn btn-danger", filter:"form-control", select:"form-control"}} onChange={this.handleDeselect} textProp="name"  valueProp="id" options={this.state.selectedRegs}  />
            </div>       
        </div>
    )
  }
}
          
ReactDOM.render(<AuthExample />, document.getElementById('app'));            