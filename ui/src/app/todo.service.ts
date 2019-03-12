import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class TodoService {

  constructor(private httpClient: HttpClient) {}

  // GET
  getTodoList() {
    return this.httpClient.get(environment.gateway + '/todo');
  }

  // POST
  addTodo(todo: Todo) {
    return this.httpClient.post(environment.gateway + '/todo', todo);
  }

  // PUT
  completeTodo(todo: Todo) {
    return this.httpClient.put(environment.gateway + '/todo', todo);
  }

  // DELETE
  deleteTodo(todo: Todo) {
    return this.httpClient.delete(environment.gateway + '/todo/' + todo.id);
  }

}

export class Todo {
  id: string;
  message: string;
  complete: boolean;
}

