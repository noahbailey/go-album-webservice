import React from 'react'; 
import ReactDOM from "react-dom";

class Table extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            albs: []
        }
    }

    componentDidMount() {
        this.refreshData();
    }

    async refreshData() {
        await fetch("http://localhost:8000/albums/")
        .then(res => res.json())
        .then(albs => this.setState({ albs }));
    }

    async deleteData(id) {
        const url = "http://localhost:8000/album/" + id
        const reqOpts = {
            method: "DELETE", 
        };
        await fetch(url, reqOpts);
        console.log("I have been deleted: #" + id)
        await this.refreshData()
    }

    editData() {
        console.log("I don't do anything yet!")
    }

    addData() {
        console.log("I don't do anything yet!")
    }

    render() {
        const tblRows = this.state.albs.map((i, row) => {
            return (
                <tr key={i.id}>
                    <td>{i.id}</td>
                    <td>{i.title}</td>
                    <td>{i.artist}</td>
                    <td>{i.price}</td>
                    <td><button class="btn btn-warning" onClick={() => this.editData(i.id)}>Edit</button></td>
                    <td><button class="btn btn-danger"onClick={() => this.deleteData(i.id)}>Delete</button></td>
                </tr>
            );
        });
        return (
            <div className="application" class="col-sm-8 col-sm-offset-2">
            <div className="table">
                <table class="table table-header">
                    <thead>
                    <tr><th>ID</th><th>Title</th><th>Artist</th><th>Price</th></tr>
                    </thead>
                    <tbody>{tblRows}</tbody>
                </table>
            </div>
            <div className="ctrls">
                <button class="btn btn-primary" onClick={() => this.refreshData()}>Refresh Data</button>
                &nbsp;
                <button class="btn btn-success" onClick={() => this.addData()}>Add Album</button>
            </div>
            </div>
        );
    }
}

ReactDOM.render(
    <Table />,
    document.getElementById('root')
);