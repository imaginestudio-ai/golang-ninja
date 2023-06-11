import * as React from "react";
import {EventCourseingForm} from "./event_booking_form";
import {Event} from "../model/event";

export interface EventCourseingFormContainerProps {
    eventID: string;
    eventServiceURL: string;
    bookingServiceURL: string;
}

export interface EventCourseingFormState {
    state: "loading"|"ready"|"saving"|"done"|"error";
    event?: Event;
}

export class EventCourseingFormContainer extends React.Component<EventCourseingFormContainerProps, EventCourseingFormState> {
    constructor(p: EventCourseingFormContainerProps) {
        super(p);

        this.state = {
            state: "loading"
        };

        fetch(p.eventServiceURL + "/events/" + p.eventID)
            .then<Event>(resp => resp.json())
            .then(event => {
                this.setState({
                    state: "ready",
                    event: event
                });
            })
    }

    render() {
        if (this.state.state === "loading") {
            return <div>Loading...</div>;
        }

        if (!this.state.event) {
            return <div>Unknown error</div>;
        }

        if (this.state.state === "done") {
            return <div className="alert alert-success">Courseing successfully completed!</div>
        }

        return <EventCourseingForm event={this.state.event} onSubmit={amount => this.handleSubmit(amount)}/>
    }

    private handleSubmit(seats: number) {
        const url = this.props.bookingServiceURL + "/events/" + this.props.eventID + "/bookings";
        const payload = {seats: seats};

        this.setState({
            event: this.state.event,
            state: "saving"
        });

        fetch(url, {method: "POST", body: JSON.stringify(payload)})
            .then(response => {
                console.log("foo")
                this.setState({
                    event: this.state.event,
                    state: response.ok ? "done" : "error"
                });
            })
    }
}