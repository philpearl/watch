

var StatusItem = React.createClass({
	render: function() {

		var status;
		if ((this.props.status == "complete") && this.props.error) {
			status = (
				<div className="row">
					<div className="col-md-12"><pre className="text-danger">{this.props.output}</pre></div>
				</div>
			)
		}
		var timestamp;
		if (this.props.status == "complete") {
			let started = new Date(this.props.started);
			let ended = new Date(this.props.ended);
			let duration = (ended - started)/1000.0;
			timestamp = <div className="col-md-1">{duration}s</div>
		}
		return (
			<div className="row">
				<div className="col-md-6">{this.props.name}</div>
				<div className="col-md-1">{this.props.type}</div>
				<div className="col-md-1">{this.props.status}</div>
				{timestamp}
				{status}
			</div>						
		)
	},
});

var StatusList = React.createClass({
	getInitialState: function() {
		return {
			status: []
		};
	},
  	render: function() {
  		var items = this.state.status.map(function(status) {
  			let key = status.type + "::" + status.name 
  			return <StatusItem key={key} {...status}></StatusItem>
  		});
    	return (
    		<div className="">
        		{items}
      		</div>
    	);
  	},
  	loadStatus: function() {
		api.status()
		.done(function(data) {
			this.setState({status: data});
		}.bind(this))
		.fail(function(jqXHR, textStatus, errorThrown) {
			console.error("status", textStatus, errorThrown.toString())
		}.bind(this));
  	},
	componentDidMount: function() {
		this.loadStatus();
    	setInterval(this.loadStatus, this.props.pollInterval);
	},  	
});
ReactDOM.render(
  <StatusList pollInterval="5000" />,
  document.getElementById('status-list')
);

