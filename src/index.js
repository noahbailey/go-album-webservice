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

    async delButton(id) {
        const url = "http://localhost:8000/album/" + id
        const reqOpts = {
            method: "DELETE", 
        };
        await fetch(url, reqOpts);
        console.log("I have been deleted: #" + id)
        await this.refreshData()
    }

    editButton(id) {
        console.log("I will be edited! #: " + id)
    }

    render() {
        const tblRows = this.state.albs.map((i, row) => {
            return (
                <tr key={i.id}>
                    <td>{i.id}</td>
                    <td>{i.title}</td>
                    <td>{i.artist}</td>
                    <td>{i.price}</td>
                    <td><button onClick={() => this.editButton(i.id)}>Edit</button></td>
                    <td><button onClick={() => this.delButton(i.id)}>Delete</button></td>
                </tr>
            );
        });
        return (
            <div className="application">
            <div className="table">
                <table>
                    <thead>
                    <tr><th>ID</th><th>Title</th><th>Artist</th><th>Price</th></tr>
                    </thead>
                    <tbody>{tblRows}</tbody>
                </table>
            </div>
            <div className="ctrls">
                <button onClick={() => this.refreshData()}>Refresh Data</button>
                <button onClick={() => this.addData()}>Add Album</button>
            </div>
            </div>
        );
    }
}

ReactDOM.render(
    <Table />,
    document.getElementById('root')
);