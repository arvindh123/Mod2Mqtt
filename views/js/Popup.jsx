class Popup extends React.Component {
    render() {
      return (
        <div class="modal-dialog modal-sm" >
          <div class="modal-content">
            <div class="model-header">
              <button class="close" onClick={this.props.closePopup}>&times;</button>
              <h4 class="modal-title  text-center">  {this.props.heading} </h4>
            </div>
            <div class="modal-body" > {this.props.text} </div>
            <div class="modal-footer"> <button class={"btn" + this.props.btnstyle} onClick={this.props.closePopup}> {this.props.btntext} </button></div>
          </div>
        </div>
  
      );
    }
  }