// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stack

import (
	"fmt"
	"time"
)

// StackElement 定义栈中的元素结构
type StackElement struct {
	Element any
	Time    time.Time
}

// Stack 定义栈结构
type Stack struct {
	elements      []*StackElement
	cleaner       chan bool
	cleanInterval time.Duration
}

// NewStack 创建一个新的栈
func NewStack(cleanInterval time.Duration) *Stack {
	s := &Stack{
		elements:      make([]*StackElement, 0),
		cleaner:       make(chan bool),
		cleanInterval: cleanInterval,
	}
	// 启动协程定期清理超时元素
	go s.startCleaner()
	return s
}

// Push 向栈中添加一个元素
func (s *Stack) Push(element any) {
	s.elements = append(s.elements, &StackElement{
		Element: element,
		Time:    time.Now().Add(s.cleanInterval),
	})
	s.cleaner <- true
}

// Pop 从栈中移除并返回最后一个添加的元素（栈顶元素）
func (s *Stack) Pop() (any, error) {
	s.cleanup()
	if s.IsEmpty() {
		return nil, fmt.Errorf("pop from an empty stack")
	}
	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return element, nil
}

// Peek 返回栈顶元素但不移除它
func (s *Stack) Peek() (any, error) {
	s.cleanup()
	if s.IsEmpty() {
		return nil, fmt.Errorf("peek from an empty stack")
	}
	return s.elements[len(s.elements)-1], nil
}

// IsEmpty 检查栈是否为空
func (s *Stack) IsEmpty() bool {
	s.cleanup()
	return len(s.elements) == 0
}

// Size 返回栈中元素的数量
func (s *Stack) Size() int {
	s.cleanup()
	return len(s.elements)
}

// cleanup 清理超时的元素
func (s *Stack) cleanup() {
	currentTime := time.Now()
	index := 0
	for _, element := range s.elements {
		if element.Time.After(currentTime) {
			s.elements[index] = element
			index++
		}
	}
	s.elements = s.elements[:index]
}

// startCleaner 启动协程定期清理超时元素
func (s *Stack) startCleaner() {
	ticker := time.NewTicker(s.cleanInterval)
	for range ticker.C {
		select {
		case <-s.cleaner:
			s.cleanup()
		default:
			s.cleanup()
		}
	}
}
