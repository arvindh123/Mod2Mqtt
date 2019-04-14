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
          <PrivateRoute path="/protected1" component={Protected1} />
          <PrivateRoute path="/protected2" component={Protected2} />
          <PrivateRoute path="/protected3" component={Protected3} />
          <PrivateRoute path="/protected4" component={Protected4} />
          <PrivateRoute path="/protected5" component={Protected5} />

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

    console.log("cWS Connnecting",this.conn)
    if (this.conn.hasOwnProperty('readyState')){
      if (this.conn.readState !==0) {
        this.conn = new  WebSocket('ws://localhost:5000/ws');

      }
    }else{
      this.conn = new  WebSocket('ws://localhost:5000/ws');
    }
    
    // conn.
    // return conn
  }, 
  WsClose(){
    this.conn.close()
  },
  WsSendCmd(cmd){
    this.conn.send( JSON.stringify({"Cmd"  : parseInt(cmd)})  )
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
              <Link to="/protected1">Settings 1</Link>
            </li>
            <li>
              <Link to="/protected2">Settings 2</Link>
            </li>
            <li>
              <Link to="/protected3">Settings 3</Link>
            </li>
            <li>
              <Link to="/protected4">Settings 4</Link>
            </li>
            <li>
              <Link to="/protected5">Settings 5</Link>
            </li>
          </ul>

          <ul class="nav navbar-nav  navbar-right">
            <li class="dropdown"><a class="dropdown-toggle" data-toggle="dropdown" href="#">Hi User</a>
              <ul class="dropdown-menu navbar-right">
                <li><a href="#" onClick={this.props.PassFunc} > Logout</a></li>
                <li><a href="/settings/view">Mqtt & Modbus Settings</a></li>
                <li><a href="/Dev/view">Devices</a></li>
                <li><a href="/modRead/view">Modbus Read Registers FC 0x03 </a></li>
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
            <Link class="navbar-brand" to="/home">Modbus to MQTT Gateway Web Interface</Link>
          </div>

          <ul class="nav navbar-nav">
            <li>
              <Link to="/protected1">Settings 1</Link>
            </li>
            <li>
              <Link to="/protected2">Settings 2</Link>
            </li>
            <li>
              <Link to="/protected3">Settings 3</Link>
            </li>
            <li>
              <Link to="/protected4">Settings 4</Link>
            </li>
            <li>
              <Link to="/protected5">Settings 5</Link>
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
    this.send2sock = this.send2sock.bind(this)
    this.send2sockLastsent = this.send2sockLastsent.bind(this)
    this.send2sockStatuslog= this.send2sockStatuslog.bind(this)
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
          statuslog :this.state.statuslog.concat([ rec.statuslog ]).slice(-5)
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
          <h3>Last Sent MQTT Message</h3>
          <pre class="box" id="mqtt_msg"> 
          {this.state.mqlastsent}
          </pre>
          <h3>Last Aquired Modbus Message</h3>
          <pre class="box" id="mqtt_msg"> 
          {this.state.modlastaquired}
          </pre>
          <div class="">
            <div class="col-lg-6">
              <table class="table"> <tbody>
                <tr>
                  <td>Mqtt Server :</td>
                  <td id="mqtt_ip"> </td>
                </tr>
                <tr>
                  <td>Mqtt Port :</td>
                  <td id="mqtt_port"></td>
                </tr>
              </tbody> </table>
            </div>
            <div class="col-lg-6">
              <table class="table"> <tbody>
                <tr>
                  <td>Modbus Server :</td>
                  <td id="mqtt_ip"> </td>
                </tr>
                <tr>
                  <td>Modbus Port :</td>
                  <td id="mqtt_port"> </td>
                </tr>
                <tr>
                  <td>
                    Status Log
                  </td>
                  <td id="log">
                  <ul>{ this.state.statuslog.slice(-5).map( (msg, idx) => <li key={'msg-' + idx }>{ msg}</li> )}</ul>;
    
                  </td>
                </tr>
              </tbody> </table>
            </div>
          </div>
          {RAuth.isAuthenticated ? 
            <div> 

              Hi Your are Logged In
            </div> 
            : null }
        </div>
        }
}
            

class  Protected1 extends React.Component {
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

class  Protected111 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
    };
  }

  componentDidMount() {
    PostGet.Get("/user/users", "GET")
           .then(response  => response.json())
           .then(json => {this.setState({ data:json})});
  };
  

  render() { 
    const { data } = this.state;
    return (
      <div>
        
      
      <table class="table" >
        <tr> <th>ID </th> <th>Name</th> <th>Password</th> <th>UserName</th>  </tr>
        <tbody>{data.msg.map(
                  function(item, i){
                    return (
                      <tr key={i}>
                      <td >{item.ID}</td><td >{item.Name}</td><td >{item.Password}</td><td >{item.UserName}</td>
                      </tr>
                    )
                  } 
                )
        }
       
        </tbody>
      </table>
      </div>
    );
  }
}

// function Protected1() {
//   return <h3>Protected1</h3>;
//           }
class  Protected2 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      messages: [],
      tosend : '',
      username :''

    };
    this.send2sock = this.send2sock.bind(this)
    this.commonChange = this.commonChange.bind(this)
    this.send2sockLastsent = this.send2sockLastsent.bind(this)
    this.send2sockStatuslog= this.send2sockStatuslog.bind(this)
    this.send2sockCommands= this.send2sockCommands.bind(this)
    
  }
  commonChange(event) {
    // console.log(this.state.tosend)
    this.setState({
      [event.target.name]: event.target.value
    });
  }
  
  // for testing purposes: sending to the echo service which will send it back back
  send2sock () { 
    console.log(this.state.tosend)
    WsClient.WsSend({"Content"  : this.state.tosend})  
  }

  send2sockLastsent () { 
    console.log(this.state.tosend)
    WsClient.WsSend({"lastsent"  : this.state.tosend})  
  }

  send2sockStatuslog () { 
    console.log(this.state.tosend)
    WsClient.WsSend({"statuslog"  : this.state.tosend})  
  }

  send2sockCommands () { 
    console.log(this.state.tosend)
    WsClient.WsSend({"Cmd"  : this.state.tosend})  
  }
  
  componentDidMount(){
  
    WsClient.conn .onmessage = evt => { 
      // add the new message to state
      console.log("In Protected 2 WS Event", evt)
    	this.setState({
      	messages : this.state.messages.concat([ JSON.parse(evt.data)])
      })
    };
    
  }

  
  render() {
    // console.log(this.state.messages)
    // slice(-5) gives us the five most recent messages
    return  <div> 
            <table class="table"> <tbody> 
            <tr>
              <td >
                  
                  <input id="send"type="submit" value="Send" onClick = {this.send2sockLastsent}/>
                  <input id="send"type="submit" value="Send" onClick = {this.send2sockStatuslog}/>
                  <input id="send"type="submit" value="Command" onClick = {this.send2sockCommands}/>
                  <input type="text"  name="tosend" onChange={this.commonChange} /> 

                
              </td>
            </tr>
            </tbody>
            </table>
            <ul>{ this.state.messages.slice(-5).map( (msg, idx) => <li key={'msg-' + idx }>{ msg.Content}</li> )}</ul>;
            </div>
    // return <ul></ul>;
  }
}

class Protected3 extends React.Component {
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
      comport :'',
      baudrate : 0,
      databits : 0,
      parity :0,
      stopbits:0,
      timeout:0,
      
    };
    this.commonChange = this.commonChange.bind(this)
    this.handleAddReg = this.handleAddReg.bind(this)
    this.getmodregs = this.getmodregs.bind(this)
    this.UpdateModParams = this.UpdateModParams.bind(this)
    this.getmodparams = this.getmodparams.bind(this)
  
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
    this.getmodparams();
  };

  getmodparams() {
    PostGet.Get("api/v1/modbus/params", "GET")
           .then(response  => response.json())
           .then(json => { //console.log(json.msg)
                            this.setState({ modparams:json.msg[0], comport:json.msg[0].comport, baudrate:json.msg[0].baudrate, databits:json.msg[0].databits, parity:json.msg[0].parity, stopbits:json.msg[0].stopbits, timeout:json.msg[0].timeout })
                         }
                );
  }
  UpdateModParams() {
    var bod = { comport : this.state.comport, baudrate : parseInt(this.state.baudrate) , databits : parseInt(this.state.databits),  parity : this.state.parity, stopbits : parseInt(this.state.stopbits), timeout : parseInt(this.state.timeout) }
    // console.log(bod)
    PostGet.Post("api/v1/modbus/params", "POST",bod)
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
  handleAddReg(){
    let bod = { name : this.state.name , unit : parseInt(this.state.unit), functcode : parseInt(this.state.functcode), register : parseInt(this.state.register), qty:parseInt(this.state.qty), datatype:parseInt(this.state.datatype), postprocess:parseInt(this.state.postprocess), tags:parseInt(this.state.tags),}
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
        return "Float Array"
      case 2:
        return "Float Array"
      case 3:
        return "Float Array"
      case 4:
        return "Float Array"
      default:
        return "I dont know"
  }
  }
  render() { 
    return (
      <div>
         <div class="row">
         <div class="col-xs-2 panel">
        <label>Com Port</label> <input class="form-control" type="text" value={this.state.comport} name="comport" onChange = {this.commonChange }/>
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
        {/* <div class="col-xs-2 panel"> */}
        <button onClick={this.UpdateModParams} className="btn btn-primary pull-bottom" >
            Update
        </button>
        {/* </div> */}
        </div>
        <br/>
        
        <label>Name &nbsp;</label><input type="text" name="name"  onChange= { this.commonChange }/> &nbsp;
        <label>Unit &nbsp;</label><input type="number" name="unit"  onChange= { this.commonChange }/> &nbsp;  
        <label>Function Code &nbsp;</label><input type="number" name="functcode" onChange= { this.commonChange }/> &nbsp;
        <label>Register&nbsp;</label><input type="number" name="register"  onChange= { this.commonChange }/> &nbsp;
        <label>qty&nbsp;</label><input type="number" name="qty"  onChange= { this.commonChange }/> &nbsp;

        <label>DataType&nbsp;</label>
        <select  name="datatype" onChange={this.commonChange}>
          <option ></option>
            <option value ="1">Int</option>
            <option value ="2">Int Array</option>
            <option value ="3">Float</option>
            <option value ="4">Float Array</option>
        </select>
        
        &nbsp;

        <label>PostProcess&nbsp;</label><input type="text" name="postprocess"  onChange= { this.commonChange }/> &nbsp;

        <label>Tags&nbsp;</label><input type="text" name="tags"  onChange= { this.commonChange }/> &nbsp;


        <button onClick={this.handleAddReg} className="btn btn-primary float-right" >
            Add Register
        </button>
      <table class="table" >
        <tr> 
          <th>ID </th> <th>Name</th> <th>Unit</th> <th>Function Code</th>
          <th>Register</th> <th>Qty</th> <th>DataType</th> <th>PostProcess</th> <th>Tags</th>
        
        </tr>
        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td> <td >{item.name}</td>
          <td >{item.unit}</td> <td >{item.functcode}</td>
          <td >{item.register}</td> <td >{item.qty}</td>
          <td >{this.getDataTypeName(item.datatype)}</td> <td >{item.postprocess}</td>
          <td >{item.tags}</td>
          <td> 
          <table class="table">
          <tr> <th>ID </th> <th>Topic</th> <th>Qos</th> <th>Retain</th> </tr>
          <tbody>{item.mqtopic.map((item, i) => (
             <tr key={i}>
            <td >{item.id}</td><td >{item.topic}</td><td >{item.qos}</td><td >{String(item.retain)}</td>
            </tr>
          )) }
          </tbody>
          </table>
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

class Protected4 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      mqparams : {},
      topic :0,
      qos : 0,
      retain : "false" ,
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
  };
  handleAddMqTopic(){
    let bod = { topic : this.state.topic, qos : parseInt(this.state.qos), retain : this.state.retain.toLowerCase() == 'true' ? true : false }
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
        <button onClick={this.handleAddMqTopic} className="btn btn-primary float-right" >
            Add Topic
        </button>
      <table class="table" >
        <tr> <th>ID </th> <th>Topic</th> <th>QoS</th> <th>Retain</th> </tr>
        <tbody>{this.state.data.map((item, i) => (        
          <tr key={i}>
          <td >{item.id}</td><td >{item.topic}</td><td >{item.qos}</td><td >{String(item.retain)}</td>
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

class Protected5 extends React.Component {
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

// function Protected3() {
//   return <h3>Protected3</h3>;
//           }
          
class Login extends React.Component {
              state = { redirectToReferrer: false };

            login = () => {
              RAuth.authenticate(() => {
                this.setState({ redirectToReferrer: true });
              });
            };
          
  render() {
              let { from } = this.props.location.state || {from: {pathname: "/" } };
    let {redirectToReferrer} = this.state;
            // console.log(redirectToReferrer)
    if (redirectToReferrer) return <Redirect to={from} />;
        
            return (
      <div>
              <p>You must log in to view the page at {from.pathname}</p>
              <button onClick={this.login}>Log in</button>
            </div>
            );
          }
        }
        
ReactDOM.render(<AuthExample />, document.getElementById('app'));
            
// export default AuthExample;