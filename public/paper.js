var Paper = React.createClass({
  render: function() {
    return (
    <div>
      <Page publication={this.props.publication}/>
    </div>
    );
  }  
});

var Page = React.createClass({
  nextData: {"page": "", num: 1},
  load: function(num) { 
    if (this.nextData.page != "") {
      console.log("setting data from nextData");
      this.setState({
        data: this.nextData,
      });      
    }
    else {
      console.log("getting data from server");
      $.ajax({
        url: "/paper/" + this.props.publication + "/page/" + num,
        dataType: 'json',
        cache: true,
        success: function(data) {
          this.setState({
            data: data
          });
        }.bind(this),
        error: function(xhr, status, err) {
          console.error(this.props.publication, status, err.toString());
        }.bind(this)
      });                
    }
  },
  preload: function(num) {
    console.log("preloading");
    $.ajax({
      url: "/paper/" + this.props.publication + "/page/" + num,
      dataType: 'json',
      cache: true,
      success: function(data) {
        this.nextData = data;
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.publication, status, err.toString());
      }.bind(this)
    });                    
  },
  getInitialState: function() {
    return {
      data : {"page": "", num: 1}, 
    }
  },
  componentDidMount: function() {
    this.load(0);
    this.preload(1);
  },
  prevPage: function() {
    var prev = this.state.data.num - 1;
    if (prev >= 0) {
      this.load(prev);
    }    
  },
  nextPage: function() {
    var next = this.state.data.num + 1;
    this.load(next);
    this.preload(next+1);
  },  
  render: function() {
    return (
    <div>
      <img className="img-responsive" src={ "data:image/png;base64," + this.state.data.page }/>
      <div className="row page">
        <div onClick={this.prevPage} onTouchEnd={this.prevPage} className="col-md-1 col-sm-1 col-xs-1 nav left"></div>      
        <div className="col-md-10 col-sm-10 col-xs-10"></div>
        <div onClick={this.nextPage} onTouchEnd={this.nextPage} className="col-md-1 col-sm-1 col-xs-1 nav right"></div>
      </div>
    </div>
    );
  }
});

React.render(
 <Paper publication={publication}/> , document.getElementById('content')
);
